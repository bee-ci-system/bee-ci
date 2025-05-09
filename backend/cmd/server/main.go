package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bee-ci/bee-ci-system/internal/common/ghservice"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/bee-ci/bee-ci-system/internal/common/middleware"
	"github.com/bee-ci/bee-ci-system/internal/data"
	"github.com/bee-ci/bee-ci-system/internal/server/api"
	"github.com/bee-ci/bee-ci-system/internal/server/webhook"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lmittmann/tint"
)

var jwtSecret = []byte("your-very-secret-key")

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	slog.SetDefault(setUpLogging())
	slog.Debug("server is starting...")

	serverURL := mustGetenv("SERVER_URL")
	port := mustGetenv("PORT")
	mainDomain := os.Getenv("MAIN_DOMAIN")
	frontendURL := mustGetenv("FRONTEND_URL")

	githubAppID := mustGetenvInt64("GITHUB_APP_ID")
	privateKeyBase64 := mustGetenv("GITHUB_APP_PRIVATE_KEY_BASE64")
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		slog.Error("error decoding GitHub App private key from base64", slog.Any("error", err))
		os.Exit(1)
	}
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		slog.Error("error parsing GitHub App RSA private key from PEM", slog.Any("error", err))
	}

	githubAppClientID := mustGetenv("GITHUB_APP_CLIENT_ID")
	githubAppWebhookSecret := mustGetenv("GITHUB_APP_WEBHOOK_SECRET")
	githubAppClientSecret := mustGetenv("GITHUB_APP_CLIENT_SECRET")

	dbHost := mustGetenv("DB_HOST")
	dbPort := mustGetenv("DB_PORT")
	dbUser := mustGetenv("DB_USER")
	dbPassword := mustGetenv("DB_PASSWORD")
	dbName := mustGetenv("DB_NAME")
	dbOpts := mustGetenv("DB_OPTS")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", dbHost, dbPort, dbUser, dbPassword, dbName, dbOpts)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		slog.Error("error connecting to Postgres database", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("connected to Postgres database", "host", dbHost, "port", dbPort, "user", dbUser, "name", dbName, "options", dbOpts)

	redisAddr := mustGetenv("REDIS_ADDRESS")
	redisPassword := mustGetenv("REDIS_PASSWORD")

	var tlsConfig *tls.Config
	if mustGetenv("REDIS_USE_TLS") == "true" {
		tlsConfig = &tls.Config{}
	}

	redisDB := redis.NewClient(&redis.Options{
		Addr:      redisAddr,
		Password:  redisPassword,
		TLSConfig: tlsConfig,
		DB:        0, // use default DB
	})

	err = redisDB.Ping(ctx).Err()
	if err != nil {
		slog.Error("error connecting to Redis database", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("connected to Redis database", "address", redisAddr)

	influxURL := mustGetenv("INFLUXDB_URL")
	influxToken := mustGetenv("INFLUXDB_TOKEN")
	influxBucket := mustGetenv("INFLUXDB_BUCKET")
	influxOrg := mustGetenv("INFLUXDB_ORG")
	influxClient := influxdb2.NewClient(influxURL, influxToken)
	_, err = influxClient.Health(ctx)
	if err != nil {
		slog.Error("error connecting to Influx database", slog.Any("error", err))
		os.Exit(1)

	} else {
		slog.Info("connected to Influx database", "url", influxURL)
	}

	buildRepo := data.NewPostgresBuildRepo(db)
	userRepo := data.NewPostgresUserRepo(db)
	repoRepo := data.NewPostgresRepoRepo(db)
	logsRepo := data.NewInfluxLogsRepo(influxClient, influxOrg, influxBucket)

	githubService := ghservice.NewGithubService(githubAppID, rsaPrivateKey, redisDB)

	webhooks, err := webhook.NewHandler(userRepo, repoRepo, buildRepo, githubService, mainDomain, frontendURL, githubAppClientID, githubAppClientSecret, githubAppWebhookSecret, jwtSecret)
	if err != nil {
		slog.Error("error creating webhook handler", slog.Any("error", err))
		os.Exit(1)
	}
	app := api.NewApp(buildRepo, logsRepo, repoRepo, userRepo, jwtSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "hello world\n\nthis is bee-ci backend server!\n\n")
		_, _ = fmt.Fprintf(w, "SERVER_URL: %s\nMAIN_DOMAIN: %s\nFRONTEND_URL: %s\n", serverURL, mainDomain, frontendURL)
	})
	mux.Handle("/webhook/", http.StripPrefix("/webhook", webhooks.Mux()))
	mux.Handle("/api/", http.StripPrefix("/api", app.Mux()))

	corsMux := middleware.WithCORS(mux)
	loggingMux := middleware.WithTrailingSlashes(middleware.WithLogger(corsMux))
	addr := fmt.Sprint("0.0.0.0:", port)
	slog.Info("server will listen and serve", "addr", addr)
	err = http.ListenAndServe(addr, loggingMux)
	if err != nil {
		slog.Error("failed while listening and serving", slog.Any("error", err))
		os.Exit(1)
	}
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

	flyProd := os.Getenv("FLY_APP_NAME") != ""
	if flyProd {
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
