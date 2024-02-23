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
export OSV_OUTPUT_FORMAT=json
#export OSV_CALL_ANALYSIS=false
export OSV_VERBOSITY=info
export OSV_OUTPUT_FILENAME=osv-scan.json

./bin/sbom-utils