package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/felixge/httpsnoop"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func WithTrailingSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = r.URL.Path + "/"
		}

		next.ServeHTTP(w, r)
	})
}

func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.With(slog.String("request_id", makeRequestID()))

		props := []any{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		}
		if r.Header.Get("X-GitHub-Event") != "" {
			props = append(props, slog.String("event", r.Header.Get("X-GitHub-Event")))
		}
		logger.Info("new request", props...)

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxLogger{}, logger)
		r = r.Clone(ctx)

		metrics := httpsnoop.CaptureMetrics(next, w, r) // this calls next.ServeHTTP
		props = append(props, slog.Int("code", metrics.Code))
		props = append(props, slog.String("duration", metrics.Duration.String()))

		logger.Info("request completed", props...)
	})
}

func WithWebhookSecret(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtain the signature from the request
		theirSignature := r.Header.Get("X-Hub-Signature-256")
		parts := strings.Split(theirSignature, "=")
		if len(parts) != 2 {
			http.Error(w, "invalid webhook signature", http.StatusForbidden)
			return
		}
		theirHexMac := parts[1]
		theirMac, err := hex.DecodeString(theirHexMac)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding webhook signature: %v", err), http.StatusBadRequest)
			return
		}

		// Calculate our own signature
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(payload)) // make body available for reading again

		hash := hmac.New(sha256.New, []byte(webhookSecret))
		hash.Write(payload)
		ourMac := hash.Sum(nil)

		// Compare signatures
		if !hmac.Equal(theirMac, ourMac) {
			http.Error(w, "webhook signature is invalid", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func WithAuthenticatedApp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := r.Context().Value(ctxLogger{}).(*slog.Logger)

		claims := jwt.MapClaims{
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(10 * time.Minute).Unix(),
			"iss": githubAppID,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenStr, err := token.SignedString(rsaPrivateKey)
		if err != nil {
			l.Error("error signing JWT", slog.Any("error", err))
			http.Error(w, "error signing JWT", http.StatusInternalServerError)
			return
		}

		appClient := http.Client{Transport: &BearerTransport{Token: tokenStr}}

		ctx := context.WithValue(r.Context(), ctxGHAppClient{}, appClient)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	})
}

func WithAuthenticatedAppInstallation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := r.Context().Value(ctxLogger{}).(*slog.Logger)

		// read request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}
		body := bytes.NewBuffer(b)
		r.Body = io.NopCloser(bytes.NewBuffer(bytes.Clone(b))) // make body available for reading again

		// decode body from JSON into a map
		decoder := json.NewDecoder(body)
		decoder.UseNumber()

		var payload map[string]interface{}
		err = decoder.Decode(&payload)
		if err != nil {
			l.Error("error reading body", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Error reading body: %v", err)
			return
		}

		// extract installation ID from the request body

		installation, ok := payload["installation"].(map[string]interface{})
		if !ok {
			l.Error("installation key not found in payload")
			next.ServeHTTP(w, r)
			return
		}

		installationIDStr := installation["id"].(json.Number).String()
		installationID, err := strconv.ParseInt(installationIDStr, 10, 64)
		if err != nil {
			l.Error("error parsing installation id", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Error parsing installation id: %v", err)
			return
		}

		// get app installation access token

		url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)
		appClient := r.Context().Value(ctxGHAppClient{}).(http.Client)
		res, err := appClient.Post(url, "application/json", nil)
		if err != nil {
			msg := "error calling endpoint " + url
			l.Error(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		// read response body and extract the access token

		// read body
		b, err = io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}
		body = bytes.NewBuffer(b)

		// decode body from JSON into a map
		decoder = json.NewDecoder(body)
		decoder.UseNumber()

		clear(payload)
		err = decoder.Decode(&payload)
		if err != nil {
			msg := fmt.Sprint("error reading response body from " + res.Request.URL.String())
			l.Error(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		appInstallationToken, ok := payload["token"].(string)
		if !ok {
			l.Error("token key not found in payload")
			next.ServeHTTP(w, r)
			return
		}

		appInstallationClient := http.Client{
			Transport: &BearerTransport{Token: appInstallationToken},
		}
		l.Debug("installation access token obtained", slog.Any("token", appInstallationToken))

		ctx := context.WithValue(r.Context(), ctxGHInstallationClient{}, appInstallationClient)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	})
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
