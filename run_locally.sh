#!/usr/bin/env bash

export PATH_TO_SBOM=examples/sboms/sample-sbom.json
export SCAN_SBOM_WITH_BOMBER=true
export BOMBER_OUTPUT_FORMAT=html
export BOMBER_DEBUG=true
export OUTPUT_DIRECTORY=build

./bin/sbom-utils