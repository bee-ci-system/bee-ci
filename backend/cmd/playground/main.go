package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v64/github"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Get installation token for a GitHub app
var installationID int64 = 56433158
var githubAppID int64 = 938460

func main() {
	privateKeyBase64 := os.Getenv("GITHUB_APP_PRIVATE_KEY_BASE64")
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		slog.Error("error decoding GitHub App private key from base64", slog.Any("error", err))
		os.Exit(1)
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		slog.Error("error parsing GitHub App RSA private key from PEM", slog.Any("error", err))
	}

	jwtString, err := generateSignedJWT(githubAppID, rsaPrivateKey)
	if err != nil {
		log.Fatalln("generate signed jwt:", err)
	}

	return

	appClient := http.Client{Transport: &bearerTransport{Token: jwtString}}
	gh := github.NewClient(&appClient)
	res, _, err := gh.Apps.CreateInstallationToken(context.Background(), installationID, nil)
	if err != nil {
		log.Fatalln("create app installation token:", err)
	}

	fmt.Println("got token:", *res.Token)
}

// generateSignedJWT does exactly what its name says.
//
// See also: [Generating a JSON Web Token JWT for a GitHub App]
//
// [Generating a JSON Web Token JWT for a GitHub App]: https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app
func generateSignedJWT(githubAppID int64, rsaPrivateKey *rsa.PrivateKey) (string, error) {
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		"iss": githubAppID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	fmt.Printf("Created JWT struct: %#v\n\n", *token)

	tokenStr, err := token.SignedString(rsaPrivateKey)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	fmt.Printf("generated signed JWT string: %s\n\n", tokenStr)

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
