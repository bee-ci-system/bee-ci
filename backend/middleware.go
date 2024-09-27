package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/bee-ci/bee-ci-system/internal/userid"

	"github.com/felixge/httpsnoop"

	l "github.com/bee-ci/bee-ci-system/internal/logger"
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
		logger.Debug("new request", props...)

		ctx := r.Context()
		ctx = l.WithLogger(ctx, logger)
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

		hash := hmac.New(sha256.New, []byte(githubAppWebhookSecret))
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

func WithJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logger, _ := l.FromContext(r.Context())

		tokenCookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "missing JWT token", http.StatusUnauthorized)
			return
		}

		// Verify the token
		token, err := verifyToken(tokenCookie.Value)
		if err != nil {
			http.Error(w, "JWT verification failed", http.StatusUnauthorized)
			return
		}

		subject, err := token.Claims.GetSubject()
		if err != nil {
			http.Error(w, "JWT verification failed (cannot retrieve subject)", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(subject, 10, 64)
		if err != nil {
			http.Error(w, "JWT verification failed (cannot parse subject)", http.StatusUnauthorized)
			return
		}

		// Print information about the verified token
		fmt.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)

		ctx := userid.WithUserID(r.Context(), userID)
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
