package merger

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sjramblings/yaml-merge/internal/progress"
	"gopkg.in/yaml.v3"
)

func MergeYAMLFiles(file1, file2, key string, pw progress.Writer) ([]byte, error) {
	pw.Start("Starting YAML Merge")
	pw.Step("Processing files:")
	pw.Step("   1. %s", file1)
	pw.Step("   2. %s", file2)
	pw.Step("   Key: %s", key)
	pw.End()

	// Validate inputs
	if err := validateInputs(file1, file2); err != nil {
		pw.Error("Input validation failed: %v", err)
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Read files
	data1, err := os.ReadFile(file1)
	if err != nil {
		pw.Error("Failed to read %s", file1)
		return nil, fmt.Errorf("error reading file1: %w", err)
	}
	pw.Success("Read %s (%d bytes)", file1, len(data1))

	data2, err := os.ReadFile(file2)
	if err != nil {
		pw.Error("Failed to read %s", file2)
		return nil, fmt.Errorf("error reading file2: %w", err)
	}
	pw.Success("Read %s (%d bytes)", file2, len(data2))

	// Parse YAML files
	var root1, root2 yaml.Node
	if err := yaml.Unmarshal(data1, &root1); err != nil {
		pw.Error("Failed to parse %s", file1)
		return nil, fmt.Errorf("error parsing file1: %w", err)
	}
	pw.Success("Parsed %s", file1)

	if err := yaml.Unmarshal(data2, &root2); err != nil {
		pw.Error("Failed to parse %s", file2)
		return nil, fmt.Errorf("error parsing file2: %w", err)
	}
	pw.Success("Parsed %s", file2)

	// Validate root nodes have content
	if len(root1.Content) == 0 {
		return nil, fmt.Errorf("error parsing file1: file is empty or invalid")
	}
	if len(root2.Content) == 0 {
		return nil, fmt.Errorf("error parsing file2: file is empty or invalid")
	}

	// Find and merge the arrays in both documents
	var array1, array2 *yaml.Node

	// Safely check root1
	if len(root1.Content) > 0 && len(root1.Content[0].Content) > 0 {
		for i := 0; i < len(root1.Content[0].Content); i += 2 {
			if i+1 < len(root1.Content[0].Content) && root1.Content[0].Content[i].Value == key {
				array1 = root1.Content[0].Content[i+1]
				break
			}
		}
	}

	// Safely check root2
	if len(root2.Content) > 0 && len(root2.Content[0].Content) > 0 {
		for i := 0; i < len(root2.Content[0].Content); i += 2 {
			if i+1 < len(root2.Content[0].Content) && root2.Content[0].Content[i].Value == key {
				array2 = root2.Content[0].Content[i+1]
				break
			}
		}
	}

	if array1 == nil || array2 == nil {
		return nil, fmt.Errorf("key '%s' not found in one or both files", key)
	}

	// Check if the key is a mapping (object) instead of a sequence
	if array1.Kind == yaml.MappingNode {
		return nil, fmt.Errorf("key '%s' in %s must be a sequence (list with '-' prefix), found object/mapping instead", key, file1)
	}
	if array2.Kind == yaml.MappingNode {
		return nil, fmt.Errorf("key '%s' in %s must be a sequence (list with '-' prefix), found object/mapping instead", key, file2)
	}

	// Print diagnostic information
	fmt.Printf("Found in %s (%d items):\n", file1, len(array1.Content))
	fmt.Printf("Found in %s (%d items):\n\n", file2, len(array2.Content))

	// Merge the arrays
	mergedContent := mergeArrays(array1, array2)

	// Create new root with only the merged key
	result := &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind:  yaml.MappingNode,
				Style: 0,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Value: key,
					},
					{
						Kind:    yaml.SequenceNode,
						Style:   0,
						Content: mergedContent,
					},
				},
			},
		},
	}

	// Marshal back to YAML with indent set to 4
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(4)
	if err := encoder.Encode(result); err != nil {
		return nil, fmt.Errorf("error marshaling result: %w", err)
	}

	return buf.Bytes(), nil
}

func mergeArrays(array1, array2 *yaml.Node) []*yaml.Node {
	seen := make(map[string]*yaml.Node)

	// Get the actual sequence content, handling both direct sequences and document nodes
	content1 := array1.Content
	content2 := array2.Content

	// If it's a document node, get the first content node
	if array1.Kind == yaml.DocumentNode && len(array1.Content) > 0 {
		content1 = array1.Content[0].Content
	}
	if array2.Kind == yaml.DocumentNode && len(array2.Content) > 0 {
		content2 = array2.Content[0].Content
	}

	// Process array1 first
	for _, item := range content1 {
		if name := getNodeName(item); name != "" {
			// Create a deep copy to avoid modifying the original
			seen[name] = &yaml.Node{
				Kind:    item.Kind,
				Style:   item.Style,
				Tag:     item.Tag,
				Value:   item.Value,
				Content: item.Content,
			}
		}
	}

	// Process array2, overwriting duplicates
	for _, item := range content2 {
		if name := getNodeName(item); name != "" {
			// Create a deep copy to avoid modifying the original
			seen[name] = &yaml.Node{
				Kind:    item.Kind,
				Style:   item.Style,
				Tag:     item.Tag,
				Value:   item.Value,
				Content: item.Content,
			}
		}
	}

	// Convert map back to array preserving order
	result := make([]*yaml.Node, 0, len(seen))

	// First add items from array1 maintaining their order
	for _, item := range content1 {
		if name := getNodeName(item); name != "" {
			if node, exists := seen[name]; exists {
				result = append(result, node)
				delete(seen, name) // Remove from seen to avoid duplicates
			}
		}
	}

	// Then add any remaining items from array2
	for _, item := range content2 {
		if name := getNodeName(item); name != "" {
			if node, exists := seen[name]; exists {
				result = append(result, node)
			}
		}
	}

	return result
}

func processArrayItems(items []*yaml.Node, seen map[string]*yaml.Node) {
	for _, item := range items {
		name := getNodeName(item)
		if name != "" {
			// Create a deep copy of the item
			nodeCopy := &yaml.Node{
				Kind:        item.Kind,
				Style:       0,
				Tag:         item.Tag,
				Value:       item.Value,
				Anchor:      item.Anchor,
				Alias:       item.Alias,
				LineComment: item.LineComment,
				HeadComment: item.HeadComment,
				FootComment: item.FootComment,
			}

			// Deep copy the content
			if item.Content != nil {
				nodeCopy.Content = make([]*yaml.Node, len(item.Content))
				for i, n := range item.Content {
					contentCopy := &yaml.Node{
						Kind:        n.Kind,
						Style:       0,
						Tag:         n.Tag,
						Value:       n.Value,
						Anchor:      n.Anchor,
						Alias:       n.Alias,
						LineComment: n.LineComment,
						HeadComment: n.HeadComment,
						FootComment: n.FootComment,
					}
					nodeCopy.Content[i] = contentCopy
				}
			}

			seen[name] = nodeCopy
		}
	}
}

func getNodeName(node *yaml.Node) string {
	// For test data that might be simple scalar nodes
	if node.Kind == yaml.ScalarNode {
		return node.Value
	}

	// Only proceed with mapping node logic if it's a mapping node
	if node.Kind != yaml.MappingNode || len(node.Content) < 2 {
		return ""
	}

	// Look for name field in mapping
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == "name" {
			return node.Content[i+1].Value
		}
	}
	return ""
}
