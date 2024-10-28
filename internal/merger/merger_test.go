package merger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Add this struct at the top of the file after imports
type mockProgress struct{}

func (m mockProgress) Start(operation string)                  {}
func (m mockProgress) Step(format string, a ...interface{})    {}
func (m mockProgress) Success(format string, a ...interface{}) {}
func (m mockProgress) Error(format string, a ...interface{})   {}
func (m mockProgress) End()                                    {}

func TestMergeYAMLFiles(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "yaml-merge-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		yaml1       string
		yaml2       string
		key         string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful merge of sequences",
			yaml1: `
vpcs:
  - name: vpc1
    cidr: 10.0.0.0/16
  - name: vpc2
    cidr: 10.1.0.0/16`,
			yaml2: `
vpcs:
  - name: vpc2
    cidr: 10.1.0.0/16
  - name: vpc3
    cidr: 10.2.0.0/16`,
			key:     "vpcs",
			wantErr: false,
		},
		{
			name: "error on non-sequence key",
			yaml1: `
vpcFlowLogs:
  trafficType: ALL
  maxAggregationInterval: 600`,
			yaml2: `
vpcFlowLogs:
  trafficType: ALL
  maxAggregationInterval: 600`,
			key:         "vpcFlowLogs",
			wantErr:     true,
			errContains: "must be a sequence",
		},
		{
			name: "key not found",
			yaml1: `
vpcs:
  - name: vpc1`,
			yaml2: `
vpcs:
  - name: vpc2`,
			key:         "nonexistentKey",
			wantErr:     true,
			errContains: "not found in one or both files",
		},
		{
			name:        "empty files",
			yaml1:       "",
			yaml2:       "",
			key:         "vpcs",
			wantErr:     true,
			errContains: "validation error: file is empty:", // Note the colon at the end
		},
		{
			name: "different structures",
			yaml1: `
services:
  - name: service1
    port: 80`,
			yaml2: `
services:
  port: 80`,
			key:         "services",
			wantErr:     true,
			errContains: "must be a sequence",
		},
		{
			name: "merge workloadAccounts",
			yaml1: `
workloadAccounts:
  - name: Network-Dev
    description: The Network Dev account
    email: dev@example.com
    organizationalUnit: Infrastructure/Infra-Dev    
  - name: Network-Prod
    description: The Network Prod account
    email: prod@example.com
    organizationalUnit: Infrastructure/Infra-Prod`,
			yaml2: `
workloadAccounts:
  - name: Pacs-Non_prod
    description: The Non Prod PACS account
    email: pacs@example.com
    organizationalUnit: HIS/HIS-Non-Prod
  - name: Network-Prod
    description: The Network Prod account
    email: prod@example.com
    organizationalUnit: Infrastructure/Infra-Prod`,
			key:     "workloadAccounts",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test files
			file1 := filepath.Join(tmpDir, "file1.yaml")
			file2 := filepath.Join(tmpDir, "file2.yaml")

			if err := os.WriteFile(file1, []byte(tt.yaml1), 0644); err != nil {
				t.Fatalf("Failed to write file1: %v", err)
			}
			if err := os.WriteFile(file2, []byte(tt.yaml2), 0644); err != nil {
				t.Fatalf("Failed to write file2: %v", err)
			}

			// Run merger
			_, err := MergeYAMLFiles(file1, file2, tt.key, &mockProgress{})

			// Check results
			if tt.wantErr {
				if err == nil {
					t.Errorf("MergeYAMLFiles() expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("MergeYAMLFiles() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("MergeYAMLFiles() unexpected error: %v", err)
			}
		})
	}
}

func TestMergeArrays(t *testing.T) {
	tests := []struct {
		name     string
		yaml1    string
		yaml2    string
		wantLen  int
		wantKeys []string
	}{
		{
			name: "merge with no duplicates",
			yaml1: `
- name: item1
  value: val1
- name: item2
  value: val2`,
			yaml2: `
- name: item3
  value: val3
- name: item4
  value: val4`,
			wantLen:  4,
			wantKeys: []string{"item1", "item2", "item3", "item4"},
		},
		{
			name: "merge with duplicates",
			yaml1: `
- name: item1
  value: val1
- name: item2
  value: val2`,
			yaml2: `
- name: item2
  value: newval2
- name: item3
  value: val3`,
			wantLen:  3,
			wantKeys: []string{"item1", "item2", "item3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var node1, node2 yaml.Node
			if err := yaml.Unmarshal([]byte(tt.yaml1), &node1); err != nil {
				t.Fatalf("Failed to parse yaml1: %v", err)
			}
			if err := yaml.Unmarshal([]byte(tt.yaml2), &node2); err != nil {
				t.Fatalf("Failed to parse yaml2: %v", err)
			}

			result := mergeArrays(&node1, &node2)

			if len(result) != tt.wantLen {
				t.Errorf("mergeArrays() got len = %v, want %v", len(result), tt.wantLen)
			}

			// Check if all expected keys are present
			for _, key := range tt.wantKeys {
				found := false
				for _, node := range result {
					if name := getNodeName(node); name == key {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("mergeArrays() missing expected key: %v", key)
				}
			}
		})
	}
}
