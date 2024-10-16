// Package ghservice exposes a simple service that makes it easy to get app installation tokens.
package ghservice

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v64/github"
)

type GithubService struct {
	logger        *slog.Logger
	httpClient    *http.Client
	githubAppID   int64
	rsaPrivateKey *rsa.PrivateKey
}

func NewGithubService(githubAppID int64, rsaPrivateKey *rsa.PrivateKey) *GithubService {
	return &GithubService{
		logger:        slog.Default(), // TODO: add some "subsystem name" to this logger
		httpClient:    &http.Client{Timeout: 10 * time.Second},
		githubAppID:   githubAppID,
		rsaPrivateKey: rsaPrivateKey,
	}
}

func (g GithubService) GetClientForInstallation(ctx context.Context, installationID int64) (*github.Client, error) {
	// token, err := getFromRedis // TODO: Get from Redis

	token, err := g.getInstallationAccessToken(ctx, installationID)
	if err != nil {
		return nil, fmt.Errorf("get installation access token: %w", err)
	}

	client := github.NewClient(&http.Client{
		Transport: &bearerTransport{Token: token},
	})

	return client, nil
}

// TODO: Cache the token in some KV store, for example Redis. Before returning
// it, always check if 1 hour has passed.

// getInstallationAccessToken returns the installation access token for the [installationID].
//
// The token returned is short-lived – per GitHub docs, it expires after 1 hour.
func (g GithubService) getInstallationAccessToken(ctx context.Context, installationID int64) (string, error) {
	jwtString, err := g.generateSignedJWT(g.githubAppID, g.rsaPrivateKey)
	if err != nil {
		return "", fmt.Errorf("generate signed jwt: %w", err)
	}

	appClient := http.Client{Transport: &bearerTransport{Token: jwtString}}
	gh := github.NewClient(&appClient)
	res, _, err := gh.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return "", fmt.Errorf("create app installation token: %w", err)
	}

	return *res.Token, nil
}

func (g GithubService) generateSignedJWT(githubAppID int64, rsaPrivateKey *rsa.PrivateKey) (string, error) {
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
