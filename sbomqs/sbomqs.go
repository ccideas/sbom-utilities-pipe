package sbomqs

import (
	"log"
	"os"
	"sbom-utilities/utils"
)

// configure logger
var (
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

func CheckSbomqsVersion() {
	err := utils.RunLiveBashCommand("sbomqs version", "")

	if err != nil {
		Log.Fatal("ERROR: sbomqs is not installed")
	}

	Log.Print("INFO: sbomqs found")
}

func GenSbomqsArgs() string {
	var sbomqsArgs string

	sbomqsArgs += "score"

	if sbomqsOutputFormat, isSet := utils.CheckEnvVar("SBOMQS_OUTPUT_FORMAT"); isSet {
		sbomqsArgs += " --" + sbomqsOutputFormat + " "
	}

	return sbomqsArgs
}

func ScanWithSbomqs(sbom string, switches string, outputFilename string, logger *log.Logger) (result bool) {
	cmd := "sbomqs " + switches + " " + sbom
	logger.Print("running the following command: " + cmd)

	err := utils.RunLiveBashCommand(cmd, outputFilename)

	return err == nil
}

func GenOutputFilename() string {

	outputfilename := "sbomqs-results"

	repoSlug := utils.GetBitbucketRepoSlug()

	if repoSlug != "" {
		outputfilename += "_" + repoSlug
	}

	currentTime := utils.GetUTCTime()
	timeString := currentTime.Format("20060102-15-04-05")

	fileExtension := ".txt"
	if sbomqsOutputFormat, isSet := utils.CheckEnvVar("SBOMQS_OUTPUT_FORMAT"); isSet {
		if sbomqsOutputFormat == "json" {
			fileExtension = ".json"
		} else {
			fileExtension = ".txt"
		}
	}

	outputfilename += "_" + timeString + fileExtension

	return outputfilename
}
