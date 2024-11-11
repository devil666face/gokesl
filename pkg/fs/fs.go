package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func WriteFile(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0777); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}

func TempPath() string {
	return filepath.Join("/", "tmp", uuid.NewString())
}

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else {
		if os.IsNotExist(err) {
			return false
		}
	}
	return false
}
