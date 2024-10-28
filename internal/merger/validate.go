package merger

import (
	"fmt"
	"os"
)

func validateInputs(file1, file2 string) error {
	// Check if files exist
	if _, err := os.Stat(file1); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", file1)
	}
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", file2)
	}
	return nil
}
