package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	LogInfo  = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogDebug = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func RunBashCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()

	if err != nil {
		LogError.Print(err)
		return "", err
	}

	outputStr := strings.TrimSpace(string(output))
	fmt.Println(outputStr)
	return outputStr, nil
}

func CheckEnvVar(env string) (string, bool) {
	value, result := os.LookupEnv(env)
	if value == "" {
		return "", result
	}
	return value, result
}

func CheckIfEnvVarIsTrue(env string) bool {
	value := os.Getenv(env)

	return value == "true"
}

func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)

	LogInfo.Print("Checking if the following file exists: " + filename)

	if os.IsNotExist(err) {
		LogInfo.Print("file " + filename + " does not exist")
	} else if err != nil {
		LogError.Print("an error occured trying to check if " + filename + " exists " + err.Error())
	} else {
		return true
	}

	return false
}

func RunLiveBashCommand(command string) error {

	cmd := exec.Command("bash", "-c", command)

	// Create a pipe to capture the command's stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Create a pipe to capture the command's stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Read and print the live output from stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			LogInfo.Print(line)
		}
	}()

	// Read and print the live output from stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			LogDebug.Print(line)
		}
	}()

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		LogError.Print(err)
		return err
	}

	return nil
}

func MoveFile(directoryPath string, dest string, contains string) bool {

	// Open the directory
	LogInfo.Print("attempting to open: ", directoryPath)
	dir, err := os.Open(directoryPath)
	if err != nil {
		LogError.Print("error opening directory: ", err)
		return false
	}
	defer dir.Close()

	// Read the directory entries
	files, err := dir.Readdir(-1) // -1 means to read all entries
	if err != nil {
		LogError.Print("error reading directory:", err)
		return false
	}

	//create the destination dir if it does not exit
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		// Directory doesn't exist, create it
		if err := os.Mkdir(dest, 0755); err != nil {
			LogError.Print(err)
			return false
		}
		LogInfo.Print("Directory created:", dest)
	} else if err != nil {
		LogError.Print(err)
		return false
	}

	// Iterate through the files
	for _, file := range files {
		if strings.Contains(file.Name(), contains) {
			sourcePath := filepath.Join(directoryPath, file.Name())
			// Construct the destination path
			destinationFile := filepath.Join(dest, file.Name())
			err := os.Rename(sourcePath, destinationFile)
			if err != nil {
				LogError.Println("error moving file:", err)
			} else {
				LogInfo.Println("file moved successfully:", file.Name())
			}
		}
	}

	return true
}

func VerifyOrCreateDirectory(directory string) bool {

	_, err := os.Stat(directory)

	if os.IsNotExist(err) {
		LogInfo.Println("directory does not exist, it will be created")
		err := os.Mkdir(directory, 0750)
		if err != nil && !os.IsExist(err) {
			LogError.Print("error occured while creating " + directory)
			return false
		}
	} else if err != nil {
		LogError.Print(err)
		return false
	} else {
		LogInfo.Println("directory already exists")
		return true
	}

	return true
}
