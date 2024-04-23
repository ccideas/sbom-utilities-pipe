package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRunBashCommand(t *testing.T) {
	//given
	command := "echo 'Hello, World!'"

	//when
	output, err := RunBashCommand(command)

	//then
	assert.NoError(t, err, "expected no error but got one")
	expectedOutput := "Hello, World!"
	assert.Equal(t, expectedOutput, output)
}

func TestRunBashCommandFailureCase(t *testing.T) {
	//given
	command := "'Hello, World!'"

	//when
	output, err := RunBashCommand(command)

	//then
	assert.Error(t, err, "expexted an error but did not get one")
	expectedOutput := ""
	assert.Equal(t, expectedOutput, output)
}

func TestCheckEnvVarThatExists(t *testing.T) {
	//given
	os.Setenv("TEST_ENV_VAR", "1")

	//when
	_, result := CheckEnvVar("TEST_ENV_VAR")

	//then
	assert.True(t, result)
}

func TestCheckEnvVarThatDoesNotExists(t *testing.T) {
	//when
	_, result := CheckEnvVar("THIS_ENV_VAR_DOES_NOT_EXIST")

	//then
	assert.False(t, result)
}

func TestEnvVarIsTrue(t *testing.T) {
	//given
	os.Setenv("DO_SOMETHING", "true")

	//when & then
	assert.True(t, CheckIfEnvVarIsTrue("DO_SOMETHING"))
}

func TestEnvVarIsNotTrue(t *testing.T) {
	//given
	os.Setenv("DO_SOMETHING", "false")

	//when & then
	assert.False(t, CheckIfEnvVarIsTrue("DO_SOMETHING"))
}

func TestCheckIfFileExists(t *testing.T) {
	//given
	tempDirName := t.TempDir()

	assert.True(t, CheckFileExists(tempDirName), "directory does not exist")
}

func TestCheckIfFileDoesNotExists(t *testing.T) {
	assert.False(t, CheckFileExists("thisDirDoesNotExist"))
}

func TestRunLiveBashCommand(t *testing.T) {
	//given
	command := "echo 'Hello, World!'"

	//when
	err := RunLiveBashCommand(command, "")

	//then
	assert.NoError(t, err, "expected no error but got one")
}

func TestRunLiveBashCommandFailure(t *testing.T) {
	//given
	command := "'Hello, World!'"

	//when
	err := RunLiveBashCommand(command, "")

	//then
	assert.Error(t, err, "expexted an error but did not get one")
}

func TestVerifyOrCreateDirectory(t *testing.T) {
	// Create a temporary directory specific to this test
	testDir := t.TempDir()

	// Test case 1: VerifyOrCreateDirectory for a non-existent directory
	result := VerifyOrCreateDirectory(testDir)
	assert.True(t, result, "expected directory to be created, but it wasn't")

	// Test case 2: VerifyOrCreateDirectory for an existing directory
	result = VerifyOrCreateDirectory(testDir)
	assert.True(t, result, "expected directory to exist, but it didn't")
}

func TestSetEnvVariable(t *testing.T) {
	// set env variable - string
	envName := "THIS_IS_MY_ENV_VARIABLE"
	envValue := "THIS_IS_MY_ENV_VARIABLE_VALUE"
	result := SetEnvVariable(envName, envValue)

	assert.Empty(t, result, "expected result to be empty")
	assert.Equal(t, os.Getenv(envName), envValue, "expected env variables to match")
}

func TestCreateDir(t *testing.T) {
	dir := t.TempDir()

	assert.True(t, CreateDir(dir), "Expected directory creation to succeed")

	_, err := os.Stat(dir)

	assert.NoError(t, err, "Expected directory '%s' to be created, but it does not exist", dir)
}

func TestCreateDirFailure(t *testing.T) {
	t.Skip("Skipping test since it fails on bitbucket due to always running as the root user")
	invalidDir := "/invalidDir"

	assert.False(t, CreateDir(invalidDir), "Expected directory creation to fail")
}

func TestOpenFile(t *testing.T) {

	file, _ := os.CreateTemp("", "example.txt")
	_, _ = file.WriteString("Temporary file for testing")
	path := file.Name()

	srcFile, ioResult := OpenFile(path)
	assert.NotNil(t, srcFile, "Expected file to be opened successfully")
	assert.True(t, ioResult, "Expected ioresult to be true")

	defer os.Remove(path)
}

func TestOpenFileFailure(t *testing.T) {

	invalidFilePath := "nonexistent_file.txt"

	srcFile, ioResult := OpenFile(invalidFilePath)
	assert.Nil(t, srcFile, "Expected file to be nil")
	assert.False(t, ioResult, "Expected ioresult to be false")
}

func TestCreateFile(t *testing.T) {

	file, _ := os.CreateTemp("", "example.txt")
	_, _ = file.WriteString("Temporary file for testing")
	path := file.Name()

	destFile, result := CreateFile(path)
	assert.NotNil(t, destFile, "Expected file to be created successfully")
	assert.True(t, result, "Expected result to be true")

	_, err := os.Stat(path)
	assert.NoError(t, err, "Expected file '%s' to be created, but it does not exist", path)

	defer os.Remove(path)
}

func TestCreateFileFailure(t *testing.T) {

	invalidFilePath := "/invalid_path/test_file.txt"
	destFile, result := CreateFile(invalidFilePath)
	assert.Nil(t, destFile, "Expected file to be nil")
	assert.False(t, result, "Expected result to be false")

}

func TestCopyFileNoDelete(t *testing.T) {

	srcFile, err := os.CreateTemp("", "source_")
	assert.NoError(t, err, "Error creating source file")
	defer srcFile.Close()
	defer os.Remove(srcFile.Name())

	destFile, err := os.CreateTemp("", "destination_")
	assert.NoError(t, err, "Error creating destination file")
	defer destFile.Close()
	defer os.Remove(destFile.Name())

	_, err = srcFile.WriteString("Temporary content for testing")
	assert.NoError(t, err, "Error writing to source file")

	result := CopyFile(destFile, srcFile, false)
	assert.True(t, result, "Expected file copying to succeed")

	destFileInfo, err := destFile.Stat()
	assert.NoError(t, err, "Error getting destination file info")
	assert.NotEqual(t, 0, destFileInfo.Size(), "Expected destination file to have non-zero size")
}

func TestCopyFileDelete(t *testing.T) {

	srcFile, err := os.CreateTemp("", "source_")
	assert.NoError(t, err, "Error creating source file")
	defer srcFile.Close()
	defer os.Remove(srcFile.Name())

	destFile, err := os.CreateTemp("", "destination_")
	assert.NoError(t, err, "Error creating destination file")
	defer destFile.Close()
	defer os.Remove(destFile.Name())

	result := CopyFile(destFile, srcFile, true)
	assert.True(t, result, "Expected file copying with source deletion to succeed")

	// Check if the source file is deleted
	_, err = os.Stat(srcFile.Name())
	assert.Error(t, err, "Expected source file to be deleted")
	assert.True(t, os.IsNotExist(err), "Expected source file to be deleted")
}
