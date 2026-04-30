FROM golang:1.26-alpine3.23 AS builder

ARG CGO_ENABLED=0

WORKDIR /tmp/gh-updater

RUN go install github.com/go-delve/delve/cmd/dlv@v1.26.3

# Copy source files necessary to download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source files required for build
COPY cmd/gh-updater/ ./cmd/gh-updater
COPY internal/ ./internal
RUN go build -gcflags="all=-N -l" -o gh-updater ./cmd/gh-updater/main.go

FROM alpine:3.23 AS runtime

COPY --from=builder /tmp/gh-updater/gh-updater /usr/local/bin/gh-updater
COPY --from=builder /go/bin/dlv /

RUN chmod -R 777 /usr/local/bin/gh-updater

CMD [ "/dlv", "--listen=:40000", "--headless=true", "--continue", "--api-version=2", "--accept-multiclient", "exec", "/usr/local/bin/gh-updater" ]

# ENTRYPOINT [ "/usr/local/bin/gh-updater" ]
