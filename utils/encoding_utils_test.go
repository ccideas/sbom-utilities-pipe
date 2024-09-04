package utils

import (
	"encoding/base64"
	"os"
	"testing"
)

// TestEncodeFile_Success tests EncodeFile for a successful encoding.
func TestEncodeFile_Success(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Write some test data to the file
	testData := []byte("hello world")
	if _, err := tmpFile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data to temp file: %v", err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test the EncodeFile function
	expected := base64.StdEncoding.EncodeToString(testData)
	encoded, err := EncodeFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("EncodeFile returned an error: %v", err)
	}
	if encoded != expected {
		t.Errorf("EncodeFile returned unexpected result. Got %s, want %s", encoded, expected)
	}
}

// TestEncodeFile_FileNotFound tests EncodeFile when the file does not exist.
func TestEncodeFile_FileNotFound(t *testing.T) {
	// Test with a non-existent file
	_, err := EncodeFile("non_existent_file.txt")
	if err == nil {
		t.Fatalf("EncodeFile did not return an error for a non-existent file")
	}
}
