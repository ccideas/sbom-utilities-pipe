package grype

import (
	"log"
	"os"
	"regexp"
	"sbom-utilities/utils"
	"strings"
)

const GRYPE_ARG_FILENAME = "--file"
const GRYPE_DEFAULT_CMD_PARAMS = "--output table"

// configure logger
var (
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

func CheckGrypeVersion() {
	err := utils.RunLiveBashCommand("grype version", "")

	if err != nil {
		Log.Fatal("ERROR: grype is not installed")
	}

	Log.Print("INFO: grype found")
}

func GenGrypeArgs() string {
	grypeArgs, result := utils.CheckEnvVar("GRYPE_ARGS")

	if !result {
		Log.Print("Info: no grype args set, using defaults")
		grypeArgs = GRYPE_DEFAULT_CMD_PARAMS
	}

	return grypeArgs
}

func GenGrypeOutputFilename(orgGrypeArgs string) (grypeArgs string, grypeOutputFile string) {

	re := regexp.MustCompile(`--file\s+([^\s]+)`)

	// if its passed in with the cmd args
	if strings.Contains(orgGrypeArgs, GRYPE_ARG_FILENAME) {
		filename := re.FindStringSubmatch(orgGrypeArgs)

		if len(filename) > 1 {
			outputFile := filename[1]
			Log.Print("filename is set via cmd switch: ", outputFile)

			return orgGrypeArgs, outputFile
		}
	}

	// if its set as a env variable
	if grypeOutputFile, isSet := utils.CheckEnvVar("GRYPE_OUTPUT_FILENAME"); isSet {

		Log.Print("filename is set via env variable: ", grypeOutputFile)
		grypeArgs = orgGrypeArgs + " " + GRYPE_ARG_FILENAME + " " + grypeOutputFile

		return grypeArgs, grypeOutputFile
	}

	// if neither are ture - gen the filename
	grypeOutputFile = utils.GenerateFilename("grype-scan", "txt")
	Log.Print("WARNING: --file switch and GRYPE_OUTPUT_FILENAME environment variables were not used, automatically settings filename")
	Log.Print("Setting file extension to .txt.")
	grypeArgs = orgGrypeArgs + " " + GRYPE_ARG_FILENAME + " " + grypeOutputFile

	return grypeArgs, grypeOutputFile
}

func ScanWithGrypeScanner(sbom string, switches string, outputFilename string, logger *log.Logger) (result bool) {
	cmd := "grype " + switches + " " + sbom
	logger.Print("running the following command: " + cmd)

	err := utils.RunLiveBashCommand(cmd, "")

	return err == nil
}
