package bomber

import (
	"sbom-utilities/utils"
	"log"
	"os"
	"os/exec"
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

	if utils.CheckEnvVar(os.Getenv("BOMBER_DEBUG")) {
		bomberArgs += "--debug" + " "
	}

	if utils.CheckEnvVar(os.Getenv("BOMBER_IGNORE_FILE")) {
		bomberArgs += "--ignore-file " + os.Getenv("BOMBER_IGNORE_FILE") + " "
	}

	if utils.CheckEnvVar(os.Getenv("BOMBER_PROVIDER")) {
		bomberArgs += "--provider " + os.Getenv("BOMBER_PROVIDER") + " "
	}

	if utils.CheckEnvVar(os.Getenv("BOMBER_PROVIDER_USERNAME")) {
		bomberArgs += "--username " + os.Getenv("BOMBER_PROVIDER_USERNAME") + " "
	}

	if utils.CheckEnvVar(os.Getenv("BOMBER_PROVIDER_TOKEN")) {
		bomberArgs += "--token " + os.Getenv("BOMBER_PROVIDER_TOKEN") + " "
	}

	if utils.CheckEnvVar(os.Getenv("BOMBER_OUTPUT_FORMAT")) {
		bomberArgs += "--output " + os.Getenv("BOMBER_OUTPUT_FORMAT") + " "
	}

	return bomberArgs
}

func ScanWithBomber(sbom string, switches string, logger *log.Logger) (result bool) {
	cmd := "bomber scan " + switches + " " + sbom
	logger.Print("running the following command: " + cmd)
	utils.RunLiveBashCommand(cmd)
	return true
}
