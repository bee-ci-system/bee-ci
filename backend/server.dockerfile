FROM golang:1.26-alpine3.23 AS builder

ARG CGO_ENABLED=0

WORKDIR /tmp/server

RUN go install github.com/go-delve/delve/cmd/dlv@v1.26.3

# Copy source files necessary to download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source files required for build
COPY cmd/server/ ./cmd/server
COPY cmd/migrate/ ./cmd/migrate
COPY internal/ ./internal
RUN go build -gcflags="all=-N -l" -o server ./cmd/server/main.go
RUN go build -gcflags="all=-N -l" -o migrate ./cmd/migrate/main.go
COPY migrations/ ./migrations

FROM alpine:3.23 AS runtime

COPY --from=builder /tmp/server/server /usr/local/bin/server
COPY --from=builder /tmp/server/migrate /usr/local/bin/migrate
COPY --from=builder /tmp/server/migrations /app/migrations
COPY --from=builder /go/bin/dlv /

ENV MIGRATIONS_PATH=/app/migrations

RUN chmod -R 777 /usr/local/bin/server /usr/local/bin/migrate

CMD [ "/dlv", "--listen=:40000", "--headless=true", "--continue", "--api-version=2", "--accept-multiclient", "exec", "/usr/local/bin/server" ]

# ENTRYPOINT [ "/usr/local/bin/server" ]
