package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/bee-ci/bee-ci-system/data"
	"github.com/bee-ci/bee-ci-system/updater"
	"github.com/bee-ci/bee-ci-system/worker"

	"github.com/jmoiron/sqlx"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lmittmann/tint"

	"github.com/lib/pq"
)

var (
	serverURL string

	githubAppID            int64
	githubAppWebhookSecret string
	rsaPrivateKey          *rsa.PrivateKey

	githubAppClientID     string
	githubAppClientSecret string
)

type (
	ctxGHInstallationClient struct{}
	ctxGHAppClient          struct{}
)

var db *sqlx.DB

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	slog.SetDefault(setUpLogging())

	serverURL = mustGetenv("SERVER_URL")
	port := mustGetenv("PORT")

	slog.Debug("server is starting", slog.String("server_url", serverURL), slog.String("port", port))

	var err error
	githubAppID = mustGetenvInt64("GITHUB_APP_ID")
	githubAppClientID = mustGetenv("GITHUB_APP_CLIENT_ID")
	githubAppWebhookSecret = mustGetenv("GITHUB_APP_WEBHOOK_SECRET")
	githubAppClientSecret = mustGetenv("GITHUB_APP_CLIENT_SECRET")
	privateKeyBase64 := mustGetenv("GITHUB_APP_PRIVATE_KEY_BASE64")
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		slog.Error("error decoding GitHub App private key from base64", slog.Any("error", err))
		os.Exit(1)
	}

	rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		slog.Error("error parsing GitHub App RSA private key from PEM", slog.Any("error", err))
	}

	dbHost := mustGetenv("DB_HOST")
	dbPort := mustGetenv("DB_PORT")
	dbUser := mustGetenv("DB_USER")
	dbPassword := mustGetenv("DB_PASSWORD")
	dbName := mustGetenv("DB_NAME")
	dbOpts := mustGetenv("DB_OPTS")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", dbHost, dbPort, dbUser, dbPassword, dbName, dbOpts)
	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		slog.Error("error connecting to Postgres database", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("connected to Postgres database", "host", dbHost, "port", dbPort, "user", dbUser, "name", dbName, "options", dbOpts)

	influxURL := mustGetenv("INFLUXDB_URL")
	influxToken := mustGetenv("INFLUXDB_TOKEN")
	influxBucket := mustGetenv("INFLUXDB_BUCKET")
	influxOrg := mustGetenv("INFLUXDB_ORG")

	influxClient := influxdb2.NewClient(influxURL, influxToken)
	_, err = influxClient.Health(ctx)
	if err != nil {
		slog.Error("error connecting to Influx database", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("connected to Influx database", "url", influxURL)

	buildRepo := data.NewPostgresBuildRepo(db)
	userRepo := data.NewPostgresUserRepo(db)
	repoRepo := data.NewPostgresRepoRepo(db)
	logsRepo := data.NewInfluxLogsRepo(influxClient, influxOrg, influxBucket)

	slog.Debug("WILL PRINT LOGS")
	logs, err := logsRepo.Get(ctx, 1)
	if err != nil {
		slog.Error("error getting logs", slog.Any("error", err))
		os.Exit(1)
	}

	for _, log := range logs {
		slog.Debug(log)
	}
	slog.Debug("DID PRINT LOGS")

	githubService := updater.NewGithubService(githubAppID, rsaPrivateKey)

	minReconnectInterval := 10 * time.Second
	maxReconnectInterval := time.Minute
	dbListener := pq.NewListener(psqlInfo, minReconnectInterval, maxReconnectInterval, nil)
	listen := updater.New(dbListener, repoRepo, userRepo, buildRepo, githubService)
	go func() {
		err := listen.Start(ctx)
		if err != nil {
			slog.Error("error starting listener", slog.Any("error", err))
			panic(err)
		}
	}()

	w := worker.New(ctx, buildRepo)
	webhooks := NewWebhookHandler(userRepo, repoRepo, w, serverURL)
	app := NewApp(buildRepo, logsRepo, repoRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "hello world\n\nthis is bee-ci backend server")
	})
	mux.Handle("/webhook/", http.StripPrefix("/webhook", webhooks.Mux()))
	mux.Handle("/api/", http.StripPrefix("/api", app.Mux()))

	loggingMux := WithTrailingSlashes(WithLogger(mux))
	addr := fmt.Sprint("0.0.0.0:", port)
	slog.Info("server will listen and serve", "addr", addr)
	err = http.ListenAndServe(addr, loggingMux)
	if err != nil {
		slog.Error("failed while listening and serving", slog.Any("error", err))
		os.Exit(1)
	}
}

func mustGetenv(varname string) string {
	value := os.Getenv(varname)
	if value == "" {
		slog.Error(varname + " env var is empty or not set")
		os.Exit(1)
	}
	return value
}

func mustGetenvInt64(varname string) int64 {
	value := mustGetenv(varname)
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		slog.Error(varname+" env var is not a valid int64", slog.Any("error", err))
		os.Exit(1)
	}
	return i
}

func setUpLogging() *slog.Logger {
	// Configure logging
	logLevel := slog.LevelDebug
	gcpProd := os.Getenv("K_SERVICE") != "" // https://cloud.google.com/run/docs/container-contract#services-env-vars
	if gcpProd {
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
	}

	flyioProd := os.Getenv("FLY_APP_NAME") != ""
	if flyioProd {
		// TODO: Remove time since it's provided by fly.io
		//  https://github.com/lmittmann/tint/issues/73
		opts := tint.Options{Level: logLevel, TimeFormat: time.TimeOnly, AddSource: true}
		handler := tint.NewHandler(os.Stdout, &opts)
		return slog.New(handler)
	}

	opts := tint.Options{Level: logLevel, TimeFormat: time.TimeOnly, AddSource: true}
	handler := tint.NewHandler(os.Stdout, &opts)
	return slog.New(handler)
}
