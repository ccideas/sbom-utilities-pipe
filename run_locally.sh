#!/usr/bin/env bash

export PATH_TO_SBOM=examples/sboms/sample-sbom.json
export SCAN_SBOM_WITH_BOMBER=false
export BOMBER_OUTPUT_FORMAT=html
export BOMBER_DEBUG=true
export OUTPUT_DIRECTORY=build
export SCAN_SBOM_WITH_SBOMQS=true
export SBOMQS_OUTPUT_FORMAT=detailed
export BITBUCKET_REPO_SLUG=testProject

./bin/sbom-utils