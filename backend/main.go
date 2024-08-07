package main

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lmittmann/tint"

	_ "github.com/lib/pq"
)

const defaultPort = "8080"

var (
	githubAppID   int64
	webhookSecret string
	rsaPrivateKey *rsa.PrivateKey
)

type (
	ctxGHInstallationClient struct{}
	ctxGHAppClient          struct{}
	ctxLogger               struct{}
)

var db *sqlx.DB

func main() {
	slog.SetDefault(setUpLogging())
	slog.Info("server is starting...")

	var err error
	githubAppID, err = strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		slog.Error("APP_ID env var not set or not a valid int64", slog.Any("error", err))
		os.Exit(1)
	}
	webhookSecret = MustGetenv("GITHUB_APP_WEBHOOK_SECRET")
	privateKeyBase64 := MustGetenv("GITHUB_APP_PRIVATE_KEY_BASE64")
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		slog.Error("error decoding GitHub App private key from base64", slog.Any("error", err))
		os.Exit(1)
	}

	rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		slog.Error("error parsing GitHub App RSA private key from PEM", slog.Any("error", err))
	}

	dbHost := MustGetenv("DB_HOST")
	dbPort := MustGetenv("DB_PORT")
	dbUser := MustGetenv("DB_USER")
	dbPassword := MustGetenv("DB_PASSWORD")
	dbName := MustGetenv("DB_NAME")
	dbOpts := "sslmode=disable"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", dbHost, dbPort, dbUser, dbPassword, dbName, dbOpts)
	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		slog.Error("error connecting to database", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("connected to database", "host", dbHost, "port", dbPort, "user", dbUser, "name", dbName, "options", dbOpts)

	port := os.Getenv("PORT")
	if port == "" {
		slog.Info(fmt.Sprintf("PORT env var not set, using default port %s", defaultPort))
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.Handle("GET /{$}", http.HandlerFunc(handleIndex))
	mux.Handle("GET /github/callback", http.HandlerFunc(handleAuthCallback))
	mux.Handle("POST /webhook",
		WithWebhookSecret(
			WithAuthenticatedApp( // provides gh_app_client
				WithAuthenticatedAppInstallation( // provides gh_installation_client
					http.HandlerFunc(handleWebhook),
				),
			),
		),
	)

	api := NewApp(db)
	mux.Handle("GET /api", api.Mux())

	loggingMux := WithLogger(mux)
	err = http.ListenAndServe(fmt.Sprint("0.0.0.0:", port), loggingMux)
	if err != nil {
		slog.Error("failed to start listening", slog.Any("error", err))
		os.Exit(1)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	l := r.Context().Value(ctxLogger{}).(*slog.Logger)
	l.Info("request received", slog.String("path", r.URL.Path))
	_, _ = fmt.Fprintln(w, "hello world")
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code query parameter", http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, "Successfully authorized! Got code %s.", code)
}

func MustGetenv(varname string) string {
	value := os.Getenv(varname)
	if value == "" {
		slog.Error(varname + " env var is empty or not set")
		os.Exit(1)
	}
	return value
}

func setUpLogging() *slog.Logger {
	// Configure logging
	logLevel := slog.LevelDebug
	prod := os.Getenv("K_SERVICE") != "" // https://cloud.google.com/run/docs/container-contract#services-env-vars
	if prod {
		// Based on https://github.com/remko/cloudrun-slog
		const LevelCritical = slog.Level(12)
		opts := &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case slog.MessageKey:
					a.Key = "message"
				case slog.SourceKey:
					a.Key = "logging.googleapis.com/sourceLocation"
				case slog.LevelKey:
					a.Key = "severity"
					level := a.Value.Any().(slog.Level)
					if level == LevelCritical {
						a.Value = slog.StringValue("CRITICAL")
					}
				}
				return a
			},
		}

		gcpHandler := slog.NewJSONHandler(os.Stderr, opts)
		return slog.New(gcpHandler)
	} else {
		opts := tint.Options{Level: logLevel, TimeFormat: time.TimeOnly, AddSource: true}
		handler := tint.NewHandler(os.Stdout, &opts)
		return slog.New(handler)
	}
}

// MakeRequestID generates a short, random hash for use as a request ID.
func makeRequestID() string {
	randomData := make([]byte, 10)
	for i := range randomData {
		randomData[i] = byte(rand.Intn(256))
	}

	strHash := fmt.Sprintf("%x", sha256.Sum256(randomData))
	return strHash[:7]
}
