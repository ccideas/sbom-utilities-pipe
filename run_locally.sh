#!/usr/bin/env bash

export PATH_TO_SBOM=examples/sboms/sample-sbom.json
export SCAN_SBOM_WITH_BOMBER=false
export BOMBER_OUTPUT_FORMAT=html
export BOMBER_DEBUG=true
export OUTPUT_DIRECTORY=build
export SCAN_SBOM_WITH_SBOMQS=false
export SBOMQS_OUTPUT_FORMAT=detailed
export BITBUCKET_REPO_SLUG=testProject
export SCAN_SBOM_WITH_OSV=true
#export OSV_ARGS=""
#export OSV_OUTPUT_FILENAME=osv-scan.json
export SCAN_SBOM_WITH_GRYPE=true
export GRYPE_ARGS="--output table --file my-wonderful-grype-output.txt --add-cpes-if-none"
#export GRYPE_OUTPUT_FILENAME="my-grype-output.json"
export SEND_SBOM_TO_DTRACK=true
export DTRACK_URL=http://localhost:8081
export DTRACK_PROJECT_ID=SOMEPROJECTID
export DTRACK_API_KEY=SOMEAPIKEY

./bin/sbom-utils