// Package ghservice exposes a simple service that makes it easy to get app installation tokens.
package ghservice

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v64/github"
)

type GithubService struct {
	logger        *slog.Logger
	httpClient    *http.Client
	redisDB       *redis.Client
	githubAppID   int64
	rsaPrivateKey *rsa.PrivateKey
}

func NewGithubService(githubAppID int64, rsaPrivateKey *rsa.PrivateKey, redisDB *redis.Client) *GithubService {
	return &GithubService{
		logger:        slog.Default(), // TODO: add some "subsystem name" to this logger
		httpClient:    &http.Client{Timeout: 10 * time.Second},
		redisDB:       redisDB,
		githubAppID:   githubAppID,
		rsaPrivateKey: rsaPrivateKey,
	}
}

func (g GithubService) GetClientForInstallation(ctx context.Context, installationID int64) (*github.Client, error) {
	token, err := g.redisDB.Get(ctx, strconv.FormatInt(installationID, 10)).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("get from redis: %w", err)
		}
	}

	if token == "" {
		g.logger.Debug("installation token not found in redis. Will get a new one.", slog.String("installationID", strconv.FormatInt(installationID, 10)))

		token, err = g.getInstallationAccessToken(ctx, installationID)
		if err != nil {
			return nil, fmt.Errorf("get installation access token: %w", err)
		}

		err = g.redisDB.Set(ctx, strconv.FormatInt(installationID, 10), token, 59*time.Minute).Err()
		if err != nil {
			return nil, fmt.Errorf("set in redis: %w", err)
		}
		g.logger.Debug("persisted installation token in redis", slog.String("installationID", strconv.FormatInt(installationID, 10)))
	}

	client := github.NewClient(&http.Client{
		Transport: &bearerTransport{Token: token},
	})

	return client, nil
}

// getInstallationAccessToken returns the installation access token for the [installationID].
//
// The token returned is short-lived – per GitHub docs, it expires after 1 hour.
func (g GithubService) getInstallationAccessToken(ctx context.Context, installationID int64) (string, error) {
	jwtString, err := generateSignedJWT(g.githubAppID, g.rsaPrivateKey)
	if err != nil {
		return "", fmt.Errorf("generate signed jwt: %w", err)
	}

	appClient := http.Client{Transport: &bearerTransport{Token: jwtString}}
	gh := github.NewClient(&appClient)
	res, _, err := gh.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return "", fmt.Errorf("create new GitHub app installation token: %w", err)
	}

	return *res.Token, nil
}

func generateSignedJWT(githubAppID int64, rsaPrivateKey *rsa.PrivateKey) (string, error) {
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		"iss": githubAppID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(rsaPrivateKey)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return tokenStr, nil
}

type bearerTransport struct {
	Token     string
	Transport http.RoundTripper
}

// RoundTrip implements http.RoundTripper interface.
func (b *bearerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	clonedRequest := r.Clone(r.Context())
	clonedRequest.Header.Set("Authorization", "Bearer "+b.Token)
	clonedRequest.Header.Set("Accept", "application/json")

	if b.Transport == nil {
		b.Transport = http.DefaultTransport
	}

	return b.Transport.RoundTrip(clonedRequest)
}
