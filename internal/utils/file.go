package utils

import "os"

func IsFileExists(filePath string) bool {
	_, err := os.Stat(filePath)

	return err == nil
}
