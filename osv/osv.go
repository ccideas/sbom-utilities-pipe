package osv

import (
	"log"
	"os"
	"sbom-utilities/utils"
)

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
	var osvArgs string

	osvArgs += "scan"

	var osvOutputFormat string
	if tmp, isSet := utils.CheckEnvVar("OSV_OUTPUT_FORMAT"); isSet {
		osvArgs += " --format " + tmp + " "
		if tmp == "json" || tmp == "sarif" {
			osvOutputFormat = "json"
		} else if tmp == "markdown" {
			osvOutputFormat = "md"
		} else if tmp == "table" {
			osvOutputFormat = "txt"
		}
	}

	osvArgs += " --output " + GenOutputFilename(osvOutputFormat)

	if osvVerbosity, isSet := utils.CheckEnvVar("OSV_VERBOSITY"); isSet {
		osvArgs += " --verbosity " + osvVerbosity + " "
	}

	if osvCallAnalysis, isSet := utils.CheckEnvVar("OSV_CALL_ANALYSIS"); isSet {
		osvArgs += " --call-analysis " + osvCallAnalysis + " "
	}

	if osvConfigFile, isSet := utils.CheckEnvVar("OSV_CONFIG_FILE"); isSet {
		osvArgs += " --config " + osvConfigFile + " "
	}

	return osvArgs
}

func GenOutputFilename(format string) string {

	if osvOutputFile, isSet := utils.CheckEnvVar("OSV_OUTPUT_FILENAME"); isSet {
		return osvOutputFile
	} else {
		return utils.GenerateFilename("osv-scan", format)
	}
}

func ScanWithOsvScanner(sbom string, switches string, outputFilename string, logger *log.Logger) (result bool) {
	cmd := "osv-scanner " + switches + " " + " --sbom " + sbom
	logger.Print("running the following command: " + cmd)

	err := utils.RunLiveBashCommand(cmd, "")

	return err == nil
}
