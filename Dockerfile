# build stage

FROM golang:1.21.7-alpine3.18 AS builder

# copy source and build
COPY . /build
WORKDIR /build

RUN go build -o bin/sbom-utils

# runtime stage

FROM golang:1.21.7-alpine3.18

ARG ARCH

ENV SBOM_UTILITIES_MODULE_HOME="/opt/sbom-utilities" \
    BASH_VERSION="5.2.15-r5" \
    BOMBER_VERSION="0.4.8" \
    OSV_SCANNER_VERSION="v1.6.2" \
    SBOMQS_VERSION="v0.0.29"
    
RUN apk --no-cache add bash=${BASH_VERSION}

ARG BOMBER_URL="https://github.com/devops-kung-fu/bomber/releases/download/v${BOMBER_VERSION}/bomber_${BOMBER_VERSION}_linux_${ARCH}.tar.gz"
ARG BOMBER_FILENAME="bomber_${BOMBER_VERSION}_linux_${ARCH}.tar.gz"
ARG SBOMQS_URL="https://github.com/interlynk-io/sbomqs/releases/download/${SBOMQS_VERSION}/sbomqs-linux-${ARCH}"
ARG SBOMQS_FILENAME="sbomqs-linux-${ARCH}"


RUN wget ${BOMBER_URL} --quiet \
    && mkdir -p /opt/bomber \
    && tar xf ${BOMBER_FILENAME} -C /opt/bomber \
    && rm ${BOMBER_FILENAME} \
    && wget ${SBOMQS_URL} --quiet \
    && mkdir /opt/sbomqs \
    && cp ${SBOMQS_FILENAME} /opt/sbomqs \
    && chmod +x /opt/sbomqs/${SBOMQS_FILENAME} \
    && ln -s /opt/sbomqs/${SBOMQS_FILENAME} /opt/sbomqs/sbomqs \
    && chmod +x /opt/sbomqs/sbomqs \
    && go install github.com/google/osv-scanner/cmd/osv-scanner@${OSV_SCANNER_VERSION}

COPY --from=builder /build/bin/sbom-utils ${SBOM_UTILITIES_MODULE_HOME}/bin/sbom-utils

SHELL ["/bin/bash", "-c"]

# Create a non-root user and group
RUN addgroup --system --gid 1002 bitbucket-group && \
  adduser --system --uid 1002 --ingroup bitbucket-group bitbucket-user

USER bitbucket-user

WORKDIR ${SBOM_UTILITIES_MODULE_HOME}

ENV PATH="${SBOM_UTILITIES_MODULE_HOME}/bin:/opt/bomber:/opt/sbomqs:${PATH}"

ENTRYPOINT ["sbom-utils"]
