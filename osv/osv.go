package osv

import (
	"log"
	"os"
	"regexp"
	"sbom-utilities/utils"
	"strings"
)

const OSV_DEFAULT_CMD_PARAMS = "scan --format table"
const OSV_ARG_FILENAME = "--output"

// configure logger
var (
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

func CheckOsvScannerVersion() {
	err := utils.RunLiveBashCommand("osv-scanner --version", "")

	if err != nil {
		Log.Fatal("ERROR: osv-scanner is not installed")
	}

	Log.Print("INFO: osv-scanner found")
}

func GenOsvArgs() string {
	osvArgs, result := utils.CheckEnvVar("OSV_ARGS")

	if !result {
		Log.Print("Info: no osv args set, using defaults")
		osvArgs = OSV_DEFAULT_CMD_PARAMS
	}

	return osvArgs
}

func ScanWithOsvScanner(sbom string, switches string, logger *log.Logger) (result bool) {
	cmd := "osv-scanner " + switches + " " + " --sbom " + sbom
	logger.Print("running the following command: " + cmd)

	err := utils.RunLiveBashCommand(cmd, "")

	return err == nil
}

func GenOsvOutputFilename(orgOsvArgs string) (osvArgs string, osvOutputFile string) {

	re := regexp.MustCompile(`--output\s+([^\s]+)`)

	// if its passed in with the cmd args
	if strings.Contains(orgOsvArgs, OSV_ARG_FILENAME) {
		filename := re.FindStringSubmatch(orgOsvArgs)

		if len(filename) > 1 {
			outputFile := filename[1]
			Log.Print("filename is set via cmd switch: ", outputFile)

			return orgOsvArgs, outputFile
		}
	}

	// if its set as a env variable
	if osvOutputFile, isSet := utils.CheckEnvVar("OSV_OUTPUT_FILENAME"); isSet {

		Log.Print("filename is set via env variable: ", osvOutputFile)
		osvArgs = orgOsvArgs + " " + OSV_ARG_FILENAME + " " + osvOutputFile

		return osvArgs, osvOutputFile
	}

	// if neither are ture - gen the filename
	osvOutputFile = utils.GenerateFilename("osv-scan", "txt")
	Log.Print("WARNING: --output switch and OSV_OUTPUT_FILENAME environment variables were not used, automatically settings filename")
	Log.Print("Setting file extension to .txt.")
	osvArgs = orgOsvArgs + " " + OSV_ARG_FILENAME + " " + osvOutputFile

	return osvArgs, osvOutputFile
}
