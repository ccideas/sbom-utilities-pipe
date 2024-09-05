package utils

import (
	"encoding/base64"
	"fmt"
	"os"
)

// EncodeFile reads the contents of the specified file and returns it as a Base64 encoded string.
// If the file cannot be read, an error is returned.
func EncodeFile(sbomFile string) (string, error) {

	fileData, err := os.ReadFile(sbomFile)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	encodedData := base64.StdEncoding.EncodeToString(fileData)

	return encodedData, nil
}
