FROM teamserverless/license-check:0.3.9 as license-check

# Build stage
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.13-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG GIT_COMMIT
ARG VERSION

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0

WORKDIR /usr/bin/

RUN apk --no-cache add git
COPY --from=license-check /license-check /usr/bin/

WORKDIR /go/src/github.com/openfaas/faas-cli
COPY . .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }

# ldflags "-s -w" strips binary
# ldflags -X injects commit version into binary
RUN /usr/bin/license-check -path ./ --verbose=false "Alex Ellis" "OpenFaaS Author(s)"

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go test $(go list ./... | grep -v /vendor/ | grep -v /template/|grep -v /build/|grep -v /sample/) -cover

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
    go build --ldflags "-s -w \
    -X github.com/openfaas/faas-cli/version.GitCommit=${GIT_COMMIT} \
    -X github.com/openfaas/faas-cli/version.Version=${VERSION} \
    -X github.com/openfaas/faas-cli/commands.Platform=${TARGETARCH}" \
    -a -installsuffix cgo -o faas-cli

# CICD stage
FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine:3.12 as root

RUN apk --no-cache add ca-certificates git

WORKDIR /home/app

COPY --from=builder /go/src/github.com/openfaas/faas-cli/faas-cli /usr/bin/

ENV PATH=$PATH:/usr/bin/

ENTRYPOINT [ "faas-cli" ]

# Release stage
FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine:3.12 as release

RUN apk --no-cache add ca-certificates git

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk add --no-cache ca-certificates

WORKDIR /home/app

COPY --from=builder /go/src/github.com/openfaas/faas-cli/faas-cli /usr/bin/
RUN chown -R app:app ./

USER app

ENV PATH=$PATH:/usr/bin/

ENTRYPOINT ["faas-cli"]
