# Contributing to YAML Merge

Thank you for your interest in contributing to YAML Merge! This document provides guidelines and steps for contributing.

## Development Process

### 1. Setting Up Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/yaml-merge.git
   ```
3. Install development tools:
   ```bash
   make install-tools
   ```

### 2. Creating a New Feature/Fix

1. Create a new branch using:
   ```bash
   make create-branch
   ```
   - Use `feature/` prefix for new features
   - Use `fix/` prefix for bug fixes
   - Use `docs/` prefix for documentation changes

2. Make your changes following our coding standards

3. Before submitting, run checks:
   ```bash
   make check-pr
   ```

### 3. Pull Request Process

1. Create a PR using:
   ```bash
   make create-pr title="Your PR Title"
   ```

2. Your PR will automatically:
   - Run tests and build checks
   - Add appropriate labels based on branch prefix
   - Generate changelog entries

### 4. PR Requirements

- [ ] Updated tests
- [ ] Documentation updated
- [ ] Changelog entry added
- [ ] PR title follows conventional commits:
  - `feat: description` for features
  - `fix: description` for bug fixes
  - `docs: description` for documentation
  - `refactor: description` for code changes
  - `test: description` for test changes

## Release Process

Releases are automated when tags are pushed. To create a new release:

- For patch releases: `make release-patch`
- For minor releases: `make release-minor`
- For major releases: `make release-major`

## Workflow Reference

Our CI/CD pipeline (referenced in `.github/workflows/main.yml`):

1. On Pull Requests:
   - Runs unit tests
   - Performs linting checks
   - Validates documentation
   - Ensures changelog entry exists

2. On Merge to Main:
   - Updates version numbers
   - Generates release notes
   - Creates GitHub release
   - Publishes to package registries

## Code Standards

1. YAML Files
   - Use 2 space indentation
   - Include comments for complex mappings
   - Follow YAML 1.2 specification

2. Go Code
   - Follow standard Go formatting (`gofmt`)
   - Add unit tests for new functionality
   - Document public functions and types
   - Keep functions focused and small

## Getting Help

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones

## License

By contributing, you agree that your contributions will be licensed under the same terms as the project license.
