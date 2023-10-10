# build stage

FROM golang:1.21.1-alpine3.18 AS builder

# TODO: You probaly do not need bash in the builder
RUN apk update && apk upgrade \
    && apk add bash \
    && rm -rf /var/cache/apk/*

# copy source and build
COPY . /build
WORKDIR /build

RUN mkdir -p /opt/sbom-utilities-pipe/bin \
    && go build -o /opt/sbom-utilities-pipe/bin/sbom-utils

# runtime stage

FROM alpine:3.18.4

# hadolint ignore=DL3018
RUN apk update && apk upgrade \
    && apk add bash \
    && rm -rf /var/cache/apk/*

COPY --from=builder /opt/sbom-utilities-pipe/bin/sbom-utils /opt/sbom-utilities-pipe/bin/sbom-utils

SHELL ["/bin/bash", "-c"]

# Create a non-root user and group
RUN addgroup --system --gid 1002 bitbucket-group && \
  adduser --system --uid 1002 --ingroup bitbucket-group bitbucket-user

USER bitbucket-user

WORKDIR /opt/sbom-utilities-pipe
ENTRYPOINT ["bin/sbom-utils"]
