package dtrack

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sbom-utilities/utils"
)

const DTRACK_API_UPLOAD_BOM = "/api/v1/bom"

type Payload struct {
	Project string `json:"project"`
	Bom     string `json:"bom"`
}

func SendToDtrack(sbomFile string) {

	err := checkDtrackVars()
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}

	dtrackProjectId, _ := utils.CheckEnvVar("DTRACK_PROJECT_ID")
	dtrackServer, _ := utils.CheckEnvVar("DTRACK_URL")
	dtrackApiKey, _ := utils.CheckEnvVar("DTRACK_API_KEY")

	encodedSbom, err := utils.EncodeFile(sbomFile)
	if err != nil {
		log.Fatalf("ERROR: failed to encode sbom %v\n", err)
	}

	payload := Payload{
		Project: dtrackProjectId,
		Bom:     encodedSbom,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marchal JSON: %v\n", err)
	}

	req, err := utils.BuildRequest("PUT", dtrackServer+DTRACK_API_UPLOAD_BOM, jsonData)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"X-Api-Key":    dtrackApiKey,
	}
	utils.AddHeaders(req, headers)

	responseBody, status, err := utils.SendRequest(req)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Printf("Response Status: %s\n", status)
	fmt.Printf("Response Body: %s\n", responseBody)
}

func checkDtrackVars() error {

	areDtrackVarsSet := true

	if _, isSet := utils.CheckEnvVar("DTRACK_URL"); !isSet {
		log.Print("Dependency Track URL is not set. Please set DTRACK_URL")
		areDtrackVarsSet = false
	}

	if _, isSet := utils.CheckEnvVar("DTRACK_API_KEY"); !isSet {
		log.Print("Dependency Track API Key is not set. Please set DTRACK_API_KEY")
		areDtrackVarsSet = false
	}

	if _, isSet := utils.CheckEnvVar("DTRACK_PROJECT_ID"); !isSet {
		log.Print("Dependency Track Project ID is not set. Please set DTRACK_PROJECT_ID")
		areDtrackVarsSet = false
	}

	if !areDtrackVarsSet {
		return errors.New("dependency track variables are not set corectelly")
	}

	return nil
}
