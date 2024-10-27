package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bee-ci/bee-ci-system/internal/common/ghservice"

	"github.com/bee-ci/bee-ci-system/internal/data"
	"github.com/bee-ci/bee-ci-system/internal/updater"

	"github.com/jmoiron/sqlx"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"github.com/lmittmann/tint"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	slog.SetDefault(setUpLogging())
	slog.Debug("gh-updater is starting")

	frontendURL := mustGetenv("FRONTEND_URL")

	githubAppID := mustGetenvInt64("GITHUB_APP_ID")
	privateKeyBase64 := mustGetenv("GITHUB_APP_PRIVATE_KEY_BASE64")
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		slog.Error("error decoding GitHub App private key from base64", slog.Any("error", err))
		os.Exit(1)
	}

	// FIXME: remove this!!!
	slog.Debug("DEBUG: ", slog.Int64("GITHUB_APP_ID", githubAppID))
	slog.Debug("DEBUG: ", slog.String("GITHUB_APP_PRIVATE_KEY_BASE64", privateKeyBase64))

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
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
	postgresDB, err := sqlx.Connect("postgres", psqlInfo)
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

	buildRepo := data.NewPostgresBuildRepo(postgresDB)
	userRepo := data.NewPostgresUserRepo(postgresDB)
	repoRepo := data.NewPostgresRepoRepo(postgresDB)

	githubService := ghservice.NewGithubService(githubAppID, rsaPrivateKey, redisDB)

	minReconnectInterval := 10 * time.Second
	maxReconnectInterval := time.Minute
	dbListener := pq.NewListener(psqlInfo, minReconnectInterval, maxReconnectInterval, nil)
	ghUpdater := updater.New(dbListener, repoRepo, userRepo, buildRepo, githubService, frontendURL)

	err = ghUpdater.Start(ctx)
	if err != nil {
		slog.Error("error while listening", slog.Any("error", err))
		panic(err)
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
