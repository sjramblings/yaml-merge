# yaml-merge

A command-line tool for merging YAML files based on a specific key, with special handling for sequence (array) merging and duplicate detection.

## Purpose

yaml-merge is designed to merge two YAML files by combining sequences (arrays) under a specified key while handling duplicates intelligently. This is particularly useful when:

- Managing infrastructure-as-code configurations
- Combining multiple environment configurations
- Merging deployment manifests
- Consolidating configuration files

## Features

- Merges YAML sequences while preserving order
- Detects and handles duplicates based on the `name` field
- Validates input files and structure
- Provides detailed progress output
- Supports quiet mode for CI/CD pipelines
- Maintains YAML formatting and indentation

## Installation

```bash
bash
go install github.com/yourusername/yaml-merge@latest
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
  - name: Network-Prod
    description: The Network Prod account
    email: <<network-Prod-account-email>>
    organizationalUnit: Infrastructure/Infra-Prod
```

file2.yaml:
```yaml
workloadAccounts:
  - name: Network-Prod
    description: The Network Prod account
    email: <<network-Prod-account-email>>
    organizationalUnit: Infrastructure/Infra-Prod
```

### Command

```bash
yaml-merge file1.yaml file2.yaml workloadAccounts
```

### Output

```yaml
workloadAccounts:
  - name: Network-Prod
    description: The Network Prod account
    email: <<network-Prod-account-email>>
    organizationalUnit: Infrastructure/Infra-Prod
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

### Merge Algorithm

The merge process follows these steps:

1. Validates input files for existence and format
2. Parses YAML content using gopkg.in/yaml.v3
3. Merges sequences under the specified key
4. Detects and handles duplicates based on the `name` field
5. Maintains YAML formatting and indentation
6. Outputs the merged result

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Your chosen license]
