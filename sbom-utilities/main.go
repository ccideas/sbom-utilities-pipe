package main

import (
	"log"
	"os"
	"sbom-utilities/bomber"
	"sbom-utilities/utils"
)

// configure logger
var (
	LogInfo  = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
)

var DEFAULT_OUTPUT_DIRECTORY string = "../output"

func main() {

	LogInfo.Print("Starting ccideas post sbom creation actions")

	// check for sbom file or a directory where sboms can be found (user controlled)
	if !utils.CheckEnvVar(os.Getenv("PATH_TO_SBOM")) {
		if !utils.CheckFileExists(os.Getenv("PATH_TO_SBOM")) {
			LogError.Fatal("verify env var PATH_TO_SBOM is set and exists")
		}
	}

	var sbomFile string = os.Getenv("PATH_TO_SBOM")

	// setup output directory path
	var outputDir string

	if utils.CheckEnvVar(os.Getenv("OUTPUT_DIRECTORY")) {
		outputDir = os.Getenv("OUTPUT_DIRECTORY")
	} else {
		utils.VerifyOrCreateDirectory(DEFAULT_OUTPUT_DIRECTORY)
		outputDir = DEFAULT_OUTPUT_DIRECTORY
	}

	if utils.CheckIfEnvVarIsTrue(os.Getenv("SCAN_SBOM_WITH_BOMBER")) {
		scanWithBomber(sbomFile, outputDir)
	}
}

func scanWithBomber(sbomFile string, outputDir string) bool {
	LogInfo.Print("scanning sbom via bomber")
	result, output := bomber.CheckBomberVersion()

	if !result {
		LogError.Print("bomber is not installed")
	}

	LogInfo.Print(output)

	bomberArgs := bomber.GenBomberArgs()
	LogInfo.Print("the following bomber args will be used: " + bomberArgs)
	bomber.ScanWithBomber(sbomFile, bomberArgs, LogInfo)

	utils.MoveFile(".", outputDir, "bomber-results")

	return true
}
