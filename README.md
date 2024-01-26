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

### Future Tools & Featurs

The following are the next set of tools/features which will be incorported into this pipe. To requrst other tooling/features or to
vote to have a specific tool/feature integreted next, [open an issue](https://github.com/ccideas/sbom-utilities-pipe/issues)

| Tool/Feature | Description |
| ------------ | ----------- |
| [anchore/grype](https://github.com/anchore/grype) | A vulnerability scanner for container images and filesystems |
| sBOM Signing | Sign the sBOM using your priviate key to prove ownership |
| Distribution API | Send your sBOM to servers such as DependencyTrack for storage and further analysis |

And more

## YAML Definition

The following is an example of a Bitbucket Pipeline which performs the following:

1. Installes dependencies for a npm project
2. Produces a sBOM via [cyclonedx-npm-pipe](https://github.com/ccideas/cyclonedx-npm-pipe)
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
          - pipe: docker://ccideas/cyclonedx-npm-pipe:1.2.1
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
        - pipe: docker://ccideas/sbom-utilities-pipe:1.1.3
          variables:
            PATH_TO_SBOM: "build/${BITBUCKET_REPO_SLUG}.json"
            SCAN_SBOM_WITH_BOMBER: 'true' # to enable a bomber scan
            BOMBER_OUTPUT_FORMAT: 'html'
            BOMBER_DEBUG: 'true'
            OUTPUT_DIRECTORY: 'build'
            SCAN_SBOM_WITH_SBOMQS: 'true' # to enable an sbomqs scan
            SBOMQS_OUTPUT_FORMAT: 'json'
        artifacts:
          - build/*

```

## Variables

| Variable                  | Usage                                                               | Options                         | Default       | Required |
| ---------------------     | -----------------------------------------------------------         | -----------                     | -------       | -------- |
| PATH_TO_SBOM              | Used to specify the name of the sbom file to further process        | <filename>                      |               | true     |
| SCAN_SBOM_WITH_BOMBER     | Used to scan the sBOM for vulnerabilities using bomber              | true, false                     | false         | false    |
| BOMBER_DEBUG              | Used to enable debug mode during bomber scan                        | true, false                     | false         | false    |
| BOMBER_IGNORE_FILE        | Used to tell bomber what CVEs to ignore                             | <path to bomber ignore file>    | none          | false    |
| BOMBER_PROVIDER           | Used to specify what vulnerability provider bomber will use         | osv, ossindex                   | osv           | false    |
| BOMBER_PROVIDER_TOKEN     | Used to specify an API token for the selected provider              | <provider apitoken>             | none          | false    |
| BOMBER_PROVIDER_USERNAME  | Used to specify an username for the selected provider               | <provider username>             | none          | false    |
| BOMBER_OUTPUT_FORMAT      | Used to specify the output format of the bomber scan                | json, html, stdout              | stdout        | false    |
| SCAN_SBOM_WITH_SBOMQS     | Used to scan the sBOM in order to generate a quality quality score  | true, false                     | false         | false    |
| SBOMQS_OUTPUT_FORMAT      | Used to specify the output format of the sbomqs scan                | detailed, json                  | detailed      | false    |
| OUTPUT_DIRECTORY          | Used to specify the directory to place all output in                | <directory name>                | build         | false    |

## Need an sBOM

This project contains some sample sBOMs which can be found in the examples/sboms directory.
To produce a sBOM for a given project you can use the following Bitbucket Pipe

[cyclonedx-npm-pipe](https://github.com/ccideas/cyclonedx-npm-pipe)

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

A big thank-you to the teams and volunteers who make these amazing tools available
