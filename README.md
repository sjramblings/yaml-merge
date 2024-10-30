# yaml-merge

A command-line tool for merging YAML files based on a specific key, with special handling for sequence (array) merging and duplicate detection.

## Purpose

yaml-merge is designed to merge two YAML files by combining sequences (arrays) under a specified key while handling duplicates intelligently. This is particularly useful when:

- Managing infrastructure-as-code configurations
- Combining multiple environment configurations
- Merging deployment manifests
- Consolidating configuration files

## Motivation

I needed a tool to merge YAML files for a project, but found that existing tools either didn't handle duplicates or didn't preserve the order of the sequences.

This tool was developed to address these specific needs mostly in relation to my use of the AWS Landing Zone Accelerator Solution - [https://github.com/aws-ia/terraform-aws-landingzone-accelerator](https://github.com/aws-ia/terraform-aws-landingzone-accelerator) which extracts AWS resource configuration from a set of YAML files.

## Features

- Merges YAML sequences while preserving order
- Detects and handles duplicates based on the `name` field
- Validates input files and structure
- Provides detailed progress output
- Supports quiet mode for CI/CD pipelines
- Maintains YAML formatting and indentation, including anchors

## Installation

```bash
go install github.com/sjramblings/yaml-merge@latest
```

## Usage

Basic usage:

```bash
yaml-merge file1.yaml file2.yaml key
```
With quiet mode:

```bash
yaml-merge file1.yaml file2.yaml key -q
```

### Arguments

- `file1.yaml`: First YAML file to merge
- `file2.yaml`: Second YAML file to merge
- `key`: The top-level key containing the sequence to merge

### Options

- `-q, --quiet`: Suppress progress output

## Example

### Input Files

file1.yaml:
```yaml
workloadAccounts:
  - name: Network-Dev
    description: The Network Dev account
    email: <<network-Dev-account-email>>
    organizationalUnit: Infrastructure/Infra-Dev    
  - name: Network-Prod
    description: The Network Prod account
    email: <<network-Prod-account-email>>
    organizationalUnit: Infrastructure/Infra-Prod
  - name: Pacs-Non_prod
    description: The Non Prod PACS account
    email: <<workload-account-email>>
    organizationalUnit: HIS/HIS-Non-Prod
  - name: Pms-Prod
    description: The PMS prod account
    email: <<workload-account-email>>
    organizationalUnit: HIS/HIS-Prod  
```

file2.yaml:
```yaml
workloadAccounts:
  - name: Network-Dev
    description: The Network Dev account
    email: <<network-Dev-account-email>>
    organizationalUnit: Infrastructure/Infra-Dev    
  - name: Network-Prod
    description: The Network Prod account
    email: <<network-Prod-account-email>>
    organizationalUnit: Infrastructure/Infra-Prod
```

### Command

```bash
yaml-merge test/fixtures/file1.yaml test/fixtures/file2.yaml workloadAccounts
```

### Output

```bash
=== Starting YAML Merge ===
→ Processing files:
→    1. test/fixtures/file1.yaml
→    2. test/fixtures/file2.yaml
→    Key: workloadAccounts

✓ Read test/fixtures/file1.yaml (7763 bytes)
✓ Read test/fixtures/file2.yaml (15824 bytes)
✓ Parsed test/fixtures/file1.yaml
✓ Parsed test/fixtures/file2.yaml
Found in test/fixtures/file1.yaml (4 items):
Found in test/fixtures/file2.yaml (2 items):

workloadAccounts:
    - name: Network-Dev
      description: The Network Dev account
      email: <<network-Dev-account-email>>
      organizationalUnit: Infrastructure/Infra-Dev
    - name: Network-Prod
      description: The Network Prod account
      email: <<network-Prod-account-email>>
      organizationalUnit: Infrastructure/Infra-Prod
    - name: Pacs-Non_prod
      description: The Non Prod PACS account
      email: <<workload-account-email>>
      organizationalUnit: HIS/HIS-Non-Prod
    - name: Pms-Prod
      description: The PMS prod account
      email: <<workload-account-email>>
      organizationalUnit: HIS/HIS-Prod
```

## Design

The tool follows a modular design with several key components:

1. **Command Processing** (`cmd` package)
   - Handles CLI argument parsing using Cobra
   - Manages flags and options
   - Coordinates the merge operation

2. **Merger Logic** (`merger` package)
   - Validates input files
   - Parses YAML content using gopkg.in/yaml.v3
   - Performs the merge operation
   - Handles duplicate detection

3. **Progress Reporting** (`progress` package)
   - Provides feedback during operations
   - Supports quiet mode for CI/CD
   - Formats output consistently

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Apache License 2.0

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
