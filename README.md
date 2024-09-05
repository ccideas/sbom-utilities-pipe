# Bitbucket Pipelines Pipe: sBOM Utilities

A Bitbucket Pipe containing a collection of open source tools to perform
various types of additionl analysis on a CycloneDX or SPDX sBOM (Software Bill of Materials).

The official copy this project is hosted on [Bitbucket](https://bitbucket.org/ccideas1/cyclonedx-npm-pipe/src/main/).
In order to reach a diverse audience a copy of the repo also exists in [GitHub](https://github.com/ccideas/sbom-utilities-pipe).
Pull Requests should be submitted to the to the Bitbucket reposiotry and changes will be kept in sync.

## Existing toolset and roadmap

The following tooling/functionally is currently available in this pipe

### Current Tools

| Tool/Feature | Description | From Version |
| ------------ | ----------- | ----------- |
| [devops-kung-fu/bomber](https://github.com/devops-kung-fu/bomber) | Scans Software Bill of Materials (SBOMs) for security vulnerabilities | [1.0.0](https://github.com/ccideas/sbom-utilities-pipe/releases) |
| [interlynk-io/sbomqs](https://github.com/interlynk-io/sbomqs) | SBOM quality score - Quality metrics for your sboms | [1.1.1](https://github.com/ccideas/sbom-utilities-pipe/releases) |
| [osv-scanner](https://github.com/google/osv-scanner)| Vulnerability scanner which uses the data provided by the [osv.dev](https://osv.dev) | [1.2.0](https://github.com/ccideas/sbom-utilities-pipe/releases) |
| [grype](https://github.com/anchore/grype)| A vulnerability scanner for container images and filesystems | [1.4.0](https://github.com/ccideas/sbom-utilities-pipe/releases) |
| [OWASP Dependency Track](https://docs.dependencytrack.org/)| Consumes and analyzes CycloneDX BOMs at high-velocity | [1.5.0](https://github.com/ccideas/sbom-utilities-pipe/releases) |

### Future Tools & Featurs

The following are the next set of tools/features which will be incorported into this pipe. To requrst other tooling/features or to
vote to have a specific tool/feature integreted next, [open an issue](https://github.com/ccideas/sbom-utilities-pipe/issues)

| Tool/Feature | Description |
| ------------ | ----------- |
| sBOM Signing | Sign the sBOM using your priviate key to prove ownership |

And more

## YAML Definition

The following is an example of a Bitbucket Pipeline which performs the following:

1. Installes dependencies for a npm project
2. Produces a sBOM via [cyclonedx-npm-pipe](https://github.com/ccideas/cyclonedx-npm-pipe) or [cyclonedx-bitbucket-pipe](https://github.com/ccideas/cyclonedx-bitbucket-pipe)
3. Uses sbom-utilities-pipe to furter process the sBOM

In the following example the sbom-utilities-pipe scans the sBOM for vulnerabilities using
devops-kung-fu/bomber then scans the sbom to generate a quality score using interlynk-io/sbomqs.
The following code snip would need to be added to
the `bitbucket-pipelines.yml` file

```yaml
pipelines:
  default:
    - step:
        name: Build and Test
        caches:
          - node
        script:
          - npm install
          - npm test
    - step:
        name: Gen CycloneDX sBom
        caches:
          - node
        script:
          # the build directory is owned by root but the pipe runs as the bitbucket-user
          # change the permission to allow the pipe to write to the build directory
          - chmod 777 build
          - pipe: docker://ccideas/cyclonedx-npm-pipe:1.5.0
            variables:
              IGNORE_NPM_ERRORS: 'true' # optional
              NPM_SHORT_PURLS: 'true' # optional
              NPM_OUTPUT_FORMAT: 'json' # optional
              NPM_PACKAGE_LOCK_ONLY: 'false' # optional
              OUTPUT_DIRECTORY: 'build' # optional # this dir should be archived by the pipeline
        artifacts:
          - build/*
  - step:
      name: Process sBOM
      script:
        # the build directory is owned by root but the pipe runs as the bitbucket-user
        # change the permission to allow the pipe to write to the build directory
        - chmod 777 build
        - pipe: docker://ccideas/sbom-utilities-pipe:1.5.0
          variables:
            PATH_TO_SBOM: "build/${BITBUCKET_REPO_SLUG}.json"
            OUTPUT_DIRECTORY: 'build'
            # bomber config
            SCAN_SBOM_WITH_BOMBER: 'true'
            BOMBER_OUTPUT_FORMAT: 'html'
            BOMBER_DEBUG: 'true'
            # sbomqs config
            SCAN_SBOM_WITH_SBOMQS: 'true'
            SBOMQS_OUTPUT_FORMAT: 'json'
            # osv config
            SCAN_SBOM_WITH_OSV: 'true'
            OSV_OUTPUT_FORMAT: 'json'
            # grype config
            SCAN_SBOM_WITH_GRYPE: 'true'
            GRYPE_ARGS: '--output table --add-cpes-if-none'
            GRYPE_OUTPUT_FILENAME: 'grype-scan-results.txt'
            # OWASP Dependency Track
            SEND_SBOM_TO_DTRACK: 'true'
            DTRACK_URL: '<<DTRACK URL: ie - http://url:port>>'
            DTRACK_PROJECT_ID: '<<DTRACK PROJECT ID>>'
            DTRACK_API_KEY: '<<DTRACK API KEY>>'
        artifacts:
          - build/*
```

## Variables

| Variable                  | Usage                                                               | Options                         | Default       | Required |
| ---------------------     | -----------------------------------------------------------         | -----------                     | -------       | -------- |
| PATH_TO_SBOM              | Used to specify the name of the sbom file to further process        | <filename>                      |               | true     |
| OUTPUT_DIRECTORY          | Used to specify the directory to place all output in                | <directory name>                | build         | false    |
| SCAN_SBOM_WITH_BOMBER     | Used to scan the sBOM for vulnerabilities using bomber              | true, false                     | false         | false    |
| BOMBER_DEBUG              | Used to enable debug mode during bomber scan                        | true, false                     | false         | false    |
| BOMBER_IGNORE_FILE        | Used to tell bomber what CVEs to ignore                             | <path to bomber ignore file>    | none          | false    |
| BOMBER_PROVIDER           | Used to specify what vulnerability provider bomber will use         | osv, ossindex                   | osv           | false    |
| BOMBER_PROVIDER_TOKEN     | Used to specify an API token for the selected provider              | <provider apitoken>             | none          | false    |
| BOMBER_PROVIDER_USERNAME  | Used to specify an username for the selected provider               | <provider username>             | none          | false    |
| BOMBER_OUTPUT_FORMAT      | Used to specify the output format of the bomber scan                | json, html, stdout              | stdout        | false    |
| SCAN_SBOM_WITH_SBOMQS     | Used to scan the sBOM in order to generate a quality quality score  | true, false                     | false         | false    |
| SBOMQS_OUTPUT_FORMAT      | Used to specify the output format of the sbomqs scan                | detailed, json                  | detailed      | false    |
| SCAN_SBOM_WITH_OSV        | Used to scan the sBOM for vulberabilities using osv scanner         | true, false                     | false         | false    |
| OSV_ARGS                  | cmd args to use when running osv-scanner                                  | see osv-scanner scan --help for full list  |               | false    |
| OSV_OUTPUT_FILENAME       | Used to specify the filename to store the osv scan output           | <filename>                      | auto-generated| false    |
| SCAN_SBOM_WITH_GRYPE      | Used to scan the sBOM for vulberabilities using the grype scanner   | true, false                     | false         | false    |
| GRYPE_ARGS                | cmd args to use when running grype                                  | see grype --help for full list  |               | false    |
| GRYPE_OUTPUT_FILENAME     | the file to write grype ouput to                                    | <filename>                      | auto-generated| false    |
| SEND_SBOM_TO_DTRACK       | Used to send the sbom to a downstream dependency track server       | true, false                     | false         | false    |
| DTRACK_URL                | The URL includeing http/https and the port number of the DTrack API is running on | <string>          | none          | true     |
| DTRACK_PROJECT_ID         | The project id to send the sbom to in dependency track              | <string>                        | none          | true     |
| DTRACK_API_KEY            | The team API key with BOM_UPLOAD permissions                        | <string>    |                   | none          | true     |

## Support for OWASP Dependency Track

As of release 1.5.0 the sbom-utilities-pipe allows you to simpily send your CycloneDX sBOM to a OWASP Dependency track server for further
analysis. The sbom-utilities-pipe uses [dependency tracks /v1/bom PUT API](https://docs.dependencytrack.org/usage/cicd/) for the request. To use this feature it is recommended you configure
the following variables as secured repository variables in your Bitbucket project configuration.

`DTRACK_URL`
`DTRACK_PROJECT_ID`
`DTRACK_API_KEY`

Then configure your bitbucket-pipelines.yml with the following

```yaml
- step:
      name: Process sBOM
      script:
        # the build directory is owned by root but the pipe runs as the bitbucket-user
        # change the permission to allow the pipe to write to the build directory
        - chmod 777 build
        - pipe: docker://ccideas/sbom-utilities-pipe:1.4.0
          variables:
            SEND_SBOM_TO_DTRACK: 'true'
            DTRACK_URL: ${DTRACK_URL}
            DTRACK_PROJECT_ID: ${DTRACK_PROJECT_ID}
            DTRACK_API_KEY: ${DTRACK_API_KEY}
```

Once the API call is successful the response ID will be logged as such

```
Response Body: {"token":"9ad9d8f9-273f-4d99-ae16-8fc89c21cd4d"}
```

## Need an sBOM

This project contains some sample sBOMs which can be found in the examples/sboms directory.
To produce a sBOM for a given project you can use any of the following Bitbucket Pipes

* [cyclonedx-npm-pipe](https://github.com/ccideas/cyclonedx-npm-pipe)
* [cyclonedx-bitbucket-pipe](https://github.com/ccideas/cyclonedx-bitbucket-pipe)
* [syft-bitbucket-pipe](https://github.com/ccideas/syft-bitbucket-pipe)

## Live Example

A working pipeline for the popular [auditjs](https://www.npmjs.com/package/auditjs)
tool has been created as an example. The pipeline in
this fork of the [auditjs](https://www.npmjs.com/package/auditjs) tool will install the required
dependencies then generate a CycloneDX sBOM containing all the ingredients which make up the
product then the sBOM will be further processed by the sbom-utilities-pipe

* [Repository Link](https://bitbucket.org/ccideas1/fork-auditjs/src/main/)
* [Link to bitbucket-pipelines.yml](https://bitbucket.org/ccideas1/fork-auditjs/src/main/bitbucket-pipelines.yml)
* [Link to pipeline](https://bitbucket.org/ccideas1/fork-auditjs/pipelines/results/4)

## Support

If you'd like help with this pipe, or you have an issue, or a feature request, [let us know](https://github.com/ccideas/sbom-utilities-pipe/issues).

If you are reporting an issue, please include:

the version of the pipe
relevant logs and error messages
steps to reproduce

## Credits

This Bitbucket pipe is a collection and integration of the following open source tools

* [bomber](https://github.com/devops-kung-fu/bomber)
* [interlynk-io/sbomqs](https://github.com/interlynk-io/sbomqs)
* [osv-scanner](https://github.com/google/osv-scanner)
* [grype](https://github.com/anchore/grype)
* [OWASP Dependency Track](https://docs.dependencytrack.org/)

A big thank-you to the teams and volunteers who make these amazing tools available
