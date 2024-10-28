package merger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validateInputs(file1, file2 string) error {
	// Check file extensions
	for _, file := range []string{file1, file2} {
		ext := strings.ToLower(filepath.Ext(file))
		if ext != ".yaml" && ext != ".yml" {
			return fmt.Errorf("file must be a YAML file (*.yaml or *.yml): %s", file)
		}
	}

	// Check if files exist and are readable
	for _, file := range []string{file1, file2} {
		info, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("file does not exist: %s", file)
			}
			return fmt.Errorf("cannot access file %s: %w", file, err)
		}
		if info.IsDir() {
			return fmt.Errorf("path is a directory, not a file: %s", file)
		}
		if info.Size() == 0 {
			return fmt.Errorf("file is empty: %s", file)
		}
	}

	return nil
}
