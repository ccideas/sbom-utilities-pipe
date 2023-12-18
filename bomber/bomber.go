package bomber

import (
	"log"
	"os/exec"
	"sbom-utilities/utils"
	"strings"
)

func CheckBomberVersion() (bool, string) {
	cmd := exec.Command("bash", "-c", "bomber --version")
	output, err := cmd.Output()

	if err != nil {
		return false, err.Error()
	}

	outputStr := strings.TrimSpace(string(output))
	return true, outputStr
}

func GenBomberArgs() string {
	var bomberArgs string

	if _, isSet := utils.CheckEnvVar("BOMBER_DEBUG"); isSet {
		bomberArgs += "--debug" + " "
	}

	if bomberIgnoreFile, isSet := utils.CheckEnvVar("BOMBER_IGNORE_FILE"); isSet {
		bomberArgs += "--ignore-file " + bomberIgnoreFile + " "
	}

	if bomberProvider, isSet := utils.CheckEnvVar("BOMBER_PROVIDER"); isSet {
		bomberArgs += "--provider " + bomberProvider + " "
	}

	if bomberProviderUsername, isSet := utils.CheckEnvVar("BOMBER_PROVIDER_USERNAME"); isSet {
		bomberArgs += "--username " + bomberProviderUsername + " "
	}

	if bomberProviderToken, isSet := utils.CheckEnvVar("BOMBER_PROVIDER_TOKEN"); isSet {
		bomberArgs += "--token " + bomberProviderToken + " "
	}

	if bomberOutputFormat, isSet := utils.CheckEnvVar("BOMBER_OUTPUT_FORMAT"); isSet {
		bomberArgs += "--output " + bomberOutputFormat + " "
	}

	return bomberArgs
}

func ScanWithBomber(sbom string, switches string, logger *log.Logger) (result bool) {
	cmd := "bomber scan " + switches + " " + sbom
	logger.Print("running the following command: " + cmd)
	err := utils.RunLiveBashCommand(cmd, "")

	return err == nil
}
