# build stage

FROM golang:1.21.1-alpine3.18 AS builder

# copy source and build
COPY . /build
WORKDIR /build

RUN go build -o bin/sbom-utils

# runtime stage

FROM golang:1.21.1-alpine3.18

# hadolint ignore=DL3018
RUN apk update && apk upgrade \
    && apk add bash \
    && rm -rf /var/cache/apk/*

ENV SBOM_UTILITIES_MODULE_HOME="/opt/sbom-utilities"

COPY --from=builder /build/bin/sbom-utils ${SBOM_UTILITIES_MODULE_HOME}/bin/sbom-utils 

SHELL ["/bin/bash", "-c"]

# Create a non-root user and group
RUN addgroup --system --gid 1002 bitbucket-group && \
  adduser --system --uid 1002 --ingroup bitbucket-group bitbucket-user

USER bitbucket-user

WORKDIR ${SBOM_UTILITIES_MODULE_HOME}

ENV PATH="${SBOM_UTILITIES_MODULE_HOME}/bin:${PATH}"

ENTRYPOINT ["sbom-utils"]
