package main

import (
	"fmt"
	"log"
	"os"
	"sbom-utilities/bomber"
	"sbom-utilities/osv"
	"sbom-utilities/sbomqs"
	"sbom-utilities/utils"
	"sbom-utilities/version"
)

// configure logger
var (
	LogInfo  = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
)

var DEFAULT_OUTPUT_DIRECTORY string = "output"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("Version: %s\n", version.GetModuleVersion())
		return
	}

	LogInfo.Print("Starting ccideas post sbom creation actions")

	// check for sbom file or a directory where sboms can be found (user controlled)
	var sbomFile string

	if sbomFilePath, isSet := utils.CheckEnvVar("PATH_TO_SBOM"); isSet {
		if utils.CheckFileExists(sbomFilePath) {
			LogInfo.Print("found sBom: " + sbomFilePath)
			sbomFile = sbomFilePath
		} else {
			LogError.Fatal("verify env var PATH_TO_SBOM is set and the file/directory exists")
		}
	}

	// setup output directory path
	var outputDir string

	if outputDirPath, isSet := utils.CheckEnvVar("OUTPUT_DIRECTORY"); isSet {
		outputDir = outputDirPath
	} else {
		utils.VerifyOrCreateDirectory(DEFAULT_OUTPUT_DIRECTORY)
		outputDir = DEFAULT_OUTPUT_DIRECTORY
	}

	// check if a BITBUCKET_REPO_SLUG exists as it will be used throughout

	// run utilities

	// bomber
	if utils.CheckIfEnvVarIsTrue("SCAN_SBOM_WITH_BOMBER") {
		scanWithBomber(sbomFile, outputDir)
	}

	// sbomqs
	if utils.CheckIfEnvVarIsTrue("SCAN_SBOM_WITH_SBOMQS") {
		scanWithSbomqs(sbomFile, outputDir)
	}

	// osv-scanner
	if utils.CheckIfEnvVarIsTrue("SCAN_SBOM_WITH_OSV") {
		scanWithOsv(sbomFile, outputDir)
	}
}

func scanWithBomber(sbomFile string, outputDir string) bool {
	LogInfo.Print("scanning sbom via bomber")
	result, output := bomber.CheckBomberVersion()

	if !result {
		LogError.Print("bomber is not installed")
		return false
	}

	LogInfo.Print(output)

	bomberArgs := bomber.GenBomberArgs()
	LogInfo.Print("the following bomber args will be used: " + bomberArgs)
	bomber.ScanWithBomber(sbomFile, bomberArgs, LogInfo)
	currentDir, err := utils.RunBashCommand("pwd")
	if err != nil {
		LogError.Print("cant get current directory")
	}

	utils.MoveFile(currentDir, outputDir, "bomber-results")

	return true
}

func scanWithSbomqs(sbomFile string, outputDir string) (result bool) {
	LogInfo.Print("scanning sbom via sbomqs")

	// set INTERLYNK_DISABLE_VERSION_CHECK to true to diable the API call to github
	// This API call often hits a rate limit issue when running on BitBucket and
	// in turn fails the pipelien.
	// Since the sbomqs package is updated on a regular bases it is safe to disable
	// this check
	res := utils.SetEnvVariable("INTERLYNK_DISABLE_VERSION_CHECK", "true")

	if res != nil {
		LogError.Print("error when setting environment variable")
	}

	sbomqs.CheckSbomqsVersion()

	sbomqsArgs := sbomqs.GenSbomqsArgs()
	LogInfo.Print("the following sbomqs args will be used: " + sbomqsArgs)

	sbomqsFilename := sbomqs.GenOutputFilename()
	LogInfo.Print("sbomqs results will be written to: " + sbomqsFilename)
	sbomqs.ScanWithSbomqs(sbomFile, sbomqsArgs, sbomqsFilename, LogInfo)

	currentDir, err := utils.RunBashCommand("pwd")
	if err != nil {
		LogError.Print("cant get current directory")
	}

	utils.MoveFile(currentDir, outputDir, "sbomqs-results")

	return true
}

func scanWithOsv(sbomFile string, outputDir string) (result bool) {
	LogInfo.Print("scanning sbom via osv-scanner")

	osv.CheckOsvScannerVersion()

	osvArgs := osv.GenOsvArgs()
	LogInfo.Print("the following osv-scanner args will be used: " + osvArgs)

	osv.ScanWithOsvScanner(sbomFile, osvArgs, "", LogInfo)

	currentDir, err := utils.RunBashCommand("pwd")
	if err != nil {
		LogError.Print("cant get current directory")
	}

	utils.MoveFile(currentDir, outputDir, "osv-scan")

	return true
}
