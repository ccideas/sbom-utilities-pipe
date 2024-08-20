# build stage

FROM golang:1.22.6-alpine3.20 AS builder

# copy source and build
COPY . /build
WORKDIR /build

RUN go build -o bin/sbom-utils

# runtime stage

FROM golang:1.22.6-alpine3.20

ARG ARCH

ENV SBOM_UTILITIES_MODULE_HOME="/opt/sbom-utilities" \
    BASH_VERSION="5.2.26-r0" \
    BOMBER_VERSION="0.4.8" \
    OSV_SCANNER_VERSION="v1.8.3" \
    SBOMQS_VERSION="v0.1.7" \
    GRYPE_VERSION="v0.79.4"
    
ARG BOMBER_URL="https://github.com/devops-kung-fu/bomber/releases/download/v${BOMBER_VERSION}/bomber_${BOMBER_VERSION}_linux_${ARCH}.tar.gz" \
    BOMBER_FILENAME="bomber_${BOMBER_VERSION}_linux_${ARCH}.tar.gz" \
    SBOMQS_URL="https://github.com/interlynk-io/sbomqs/releases/download/${SBOMQS_VERSION}/sbomqs-linux-${ARCH}" \
    SBOMQS_FILENAME="sbomqs-linux-${ARCH}" \
    GRYPE_URL="https://raw.githubusercontent.com/anchore/grype/main/install.sh"

RUN apk update upgrade \
    && apk --no-cache add bash="${BASH_VERSION}" \
    && wget "${BOMBER_URL}" --quiet \
    && mkdir -p /opt/bomber \
    && tar xf "${BOMBER_FILENAME}" -C /opt/bomber \
    && rm "${BOMBER_FILENAME}" \
    && wget "${SBOMQS_URL}" --quiet \
    && mkdir /opt/sbomqs \
    && cp "${SBOMQS_FILENAME}" /opt/sbomqs \
    && chmod +x /opt/sbomqs/"${SBOMQS_FILENAME}" \
    && ln -s /opt/sbomqs/"${SBOMQS_FILENAME}" /opt/sbomqs/sbomqs \
    && chmod +x /opt/sbomqs/sbomqs \
    && mkdir /opt/grype \
    && wget -P /opt/grype "${GRYPE_URL}" --quiet \
    && chmod +x /opt/grype/install.sh \
    && /opt/grype/install.sh -b /opt/grype "${GRYPE_VERSION}" \
    && rm /opt/grype/install.sh \
    && go install github.com/google/osv-scanner/cmd/osv-scanner@"${OSV_SCANNER_VERSION}" \
    && go clean -cache -testcache -modcache -fuzzcache

COPY --from=builder /build/bin/sbom-utils ${SBOM_UTILITIES_MODULE_HOME}/bin/sbom-utils

SHELL ["/bin/bash", "-c"]

# Create a non-root user and group
RUN addgroup --system --gid 1002 bitbucket-group && \
  adduser --system --uid 1002 --ingroup bitbucket-group bitbucket-user

USER bitbucket-user

WORKDIR ${SBOM_UTILITIES_MODULE_HOME}

ENV PATH="${SBOM_UTILITIES_MODULE_HOME}/bin:/opt/bomber:/opt/sbomqs:/opt/grype:${PATH}"

ENTRYPOINT ["sbom-utils"]
