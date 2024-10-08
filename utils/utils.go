package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	LogInfo                       = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogError                      = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogDebug                      = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	DEFAULT_DIRECTORY_PERMISSIONS = 0755
)

func RunBashCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()

	if err != nil {
		LogError.Print(err)
		return "", err
	}

	outputStr := strings.TrimSpace(string(output))
	return outputStr, nil
}

func CheckEnvVar(env string) (string, bool) {
	value, result := os.LookupEnv(env)

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

func RunLiveBashCommand(command string, outputPath string) error {
	var outputFile io.Writer

	if outputPath != "" {
		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()
		outputFile = file
		LogInfo.Print("output will be written to " + outputPath)
	} else {
		outputFile = nil
	}

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
			if outputFile == nil {
				LogInfo.Print(line)
			} else {
				_, _ = fmt.Fprintln(outputFile, line)
			}
		}
	}()

	// Read and print the live output from stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			if outputFile == nil {
				LogDebug.Print(line)
			} else {
				_, _ = fmt.Fprintln(outputFile, line)
			}
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

func MoveFileToDestination(currentDir string, destDir string, fileToMove string, deleteSource bool) (result bool) {

	CreateDir(destDir)

	srcPath := filepath.Join(currentDir, fileToMove)
	destPath := filepath.Join(destDir, fileToMove)

	readFile, _ := OpenFile(srcPath)
	writeFile, _ := CreateFile(destPath)
	CopyFile(writeFile, readFile, deleteSource)

	defer readFile.Close()
	defer writeFile.Close()
	return result
}

func CreateDir(dir string) (result bool) {

	err := os.MkdirAll(dir, fs.FileMode(DEFAULT_DIRECTORY_PERMISSIONS))

	if err != nil {
		LogError.Print("Error creating destination directory", err)
		return false
	}

	return true
}

func OpenFile(file string) (srcFile *os.File, ioresult bool) {

	srcFile, err := os.Open(file)

	if err != nil {
		LogError.Print("Error opening file", err)
		return nil, false
	}

	return srcFile, true
}

func CreateFile(path string) (destFile *os.File, result bool) {

	destFile, err := os.Create(path)

	if err != nil {
		LogError.Print("Error creating destination file", err)
		return nil, false
	}

	//defer destFile.Close()
	return destFile, true
}

func CopyFile(dest *os.File, src *os.File, deleteSource bool) (result bool) {

	_, err := io.Copy(dest, src)

	if err != nil {
		LogError.Print("Error copying file", err)
		return false
	}

	if deleteSource {
		err = os.Remove(src.Name())

		if err != nil {
			LogError.Print("Error removing file", err)
			return false
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

func GetUTCTime() time.Time {

	return time.Now().UTC()
}

func GetBitbucketRepoSlug() string {

	if bitbucketRepoSlug, isSet := CheckEnvVar("BITBUCKET_REPO_SLUG"); isSet {
		return bitbucketRepoSlug
	} else {
		return ""
	}
}

func SetEnvVariable(name string, value string) error {

	err := os.Setenv(name, value)

	if err != nil {
		LogError.Print("Error setting environment variable: " + name)
		LogError.Print(err)
		return err
	}

	return nil
}

func GenerateFilename(prefix string, format string) string {
	repoSlug := GetBitbucketRepoSlug()

	var name string
	if repoSlug != "" {
		name = prefix + "_" + repoSlug
	} else {
		name = prefix
	}

	currentTime := GetUTCTime()
	timeString := currentTime.Format("20060102-15-04-05")

	outputfilename := name + "_" + timeString + "." + format

	return outputfilename
}
