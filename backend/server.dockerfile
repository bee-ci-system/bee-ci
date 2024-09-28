FROM golang:1.23-alpine3.20 AS builder

ARG CGO_ENABLED=0

WORKDIR /tmp/server

RUN go install github.com/go-delve/delve/cmd/dlv@v1.23.0

# Copy source files necessary to download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source files required for build
COPY *.go ./
COPY cmd/server/ ./cmd/server
COPY data/ ./data
COPY worker/ ./worker
# COPY updater/ ./updater
COPY server ./server
COPY internal/ ./internal
RUN go build -gcflags="all=-N -l" -o server ./cmd/server/main.go

FROM alpine:3.20 AS runtime

COPY --from=builder /tmp/server/server /usr/local/bin/server
COPY --from=builder /go/bin/dlv /

RUN chmod -R 777 /usr/local/bin/server

CMD [ "/dlv", "--listen=:40000", "--headless=true", "--continue", "--api-version=2", "--accept-multiclient", "exec", "/usr/local/bin/server" ]

# ENTRYPOINT [ "/usr/local/bin/server" ]
