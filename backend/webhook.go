package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/bartekpacia/ghapp/data"
	l "github.com/bartekpacia/ghapp/internal/logger"
	"github.com/bartekpacia/ghapp/worker"
)

// Define your secret key (should be stored securely, e.g., in env variables)
var jwtSecret = []byte("your-very-secret-key")

// Function to create JWT tokens with claims
func createToken(userID int64) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"iss": "bee-ci",
		// "exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(), // Issued at
	})

	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)

	tokenString, err := claims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Function to verify JWT tokens
func verifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing JWT: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid JWT")
	}

	return token, nil
}

type WebhookHandler struct {
	userRepo   data.UserRepo
	repoRepo   data.RepoRepo
	httpClient *http.Client
	worker     *worker.Worker
}

func NewWebhookHandler(userRepo data.UserRepo, repoRepo data.RepoRepo, w *worker.Worker) *WebhookHandler {
	return &WebhookHandler{
		userRepo: userRepo,
		repoRepo: repoRepo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		worker: w,
	}
}

func (h WebhookHandler) Mux() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", http.HandlerFunc(handleIndex))
	mux.Handle("GET /github/callback/", http.HandlerFunc(h.handleAuthCallback))

	mux.Handle("POST /{$}",
		WithWebhookSecret(
			WithAuthenticatedApp( // MAYBE provides gh_app_client
				WithAuthenticatedAppInstallation( // MAYBE provides gh_installation_client
					http.HandlerFunc(h.handleWebhook),
				),
			),
		),
	)

	return mux
}

func (h WebhookHandler) githubUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	const url = "https://api.github.com/user"

	logger, _ := l.FromContext(ctx)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Do the request
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	userData := map[string]interface{}{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&userData)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	// TODO: store user data in a database
	logger.Info("authenticated user",
		slog.Any("login", userData["login"]),
		slog.Any("id", userData["id"]),
		slog.Any("name", userData["name"]),
	)

	return userData, nil
}

func (h WebhookHandler) exchangeCode(ctx context.Context, code string) (userAccessToken string, err error) {
	const url = "https://github.com/login/oauth/access_token"

	reqBody := map[string]interface{}{
		"client_id":     githubAppClientID,
		"client_secret": githubAppClientSecret,
		"code":          code,
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshalling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return "", fmt.Errorf("creating new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", fmt.Errorf("unmarshalling response body: %w", err)
	}

	accessToken, ok := respBody["access_token"].(string)
	if !ok {
		errorID, _ := respBody["error"].(string)
		errorDesc, _ := respBody["error_description"].(string)

		if errorID != "" && errorDesc != "" {
			return "", fmt.Errorf("exchanging code for access token: %s: %s", errorID, errorDesc)
		}

		return "", fmt.Errorf("access token is missing or invalid")
	}

	return accessToken, nil
}

// HandleAuthCallback exercises the [web application flow] for authorizing GitHub Apps.
//
// [web application flow]: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#web-application-flow
func (h WebhookHandler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, _ := l.FromContext(ctx)

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code query parameter", http.StatusBadRequest)
		return
	}

	accessToken, err := h.exchangeCode(r.Context(), code)
	if err != nil {
		logger.Error("error exchanging code for access token", slog.Any("error", err))
		http.Error(w, "error exchanging code for access token", http.StatusInternalServerError)
		return
	}

	// Make user request to GitHub to get data
	userData, err := h.githubUserInfo(r.Context(), accessToken)
	if err != nil {
		logger.Error("error getting user info", slog.Any("error", err))
		http.Error(w, "error getting user info", http.StatusInternalServerError)
		return
	}

	// Extract user information
	userID, err := userData["id"].(json.Number).Int64()
	if err != nil {
		logger.Error("error parsing user ID", slog.Any("error", err))
		http.Error(w, "error parsing user information", http.StatusInternalServerError)
		return
	}

	username, ok := userData["login"].(string)
	if !ok {
		logger.Error("username not found or not a string")
		http.Error(w, "error parsing user information", http.StatusInternalServerError)
		return
	}

	err = h.userRepo.Upsert(ctx, data.NewUser{
		ID:           userID,
		Username:     username,
		AccessToken:  accessToken,
		RefreshToken: "", // GitHub doesn't provide refresh tokens for OAuth Apps
	})
	if err != nil {
		logger.Error("error upserting user to database", slog.Any("error", err))
		http.Error(w, "error upserting user to database", http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Successfully authorized! User %s (ID: %d) has been saved to the database.", username, userID)
	logger.Info(msg, "access_token", accessToken)

	// Create JWT
	token, err := createToken(userID)
	if err != nil {
		logger.Error("error creating token", slog.Any("error", err))
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	jwtTokenCookie := &http.Cookie{
		Name:   "jwt",
		Value:  token,
		Domain: "bee-ci.pacia.tech",
		Path:   "/",
	}

	http.SetCookie(w, jwtTokenCookie)

	http.Redirect(w, r, "https://app.bee-ci.pacia.tech/dashboard", http.StatusSeeOther)
}

func (h WebhookHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	var payload map[string]interface{}
	err := decoder.Decode(&payload)
	if err != nil {
		msg := "error reading request body"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	installationIDStr := payload["installation"].(map[string]interface{})["id"].(json.Number).String()
	installationID, err := strconv.ParseInt(installationIDStr, 10, 64)
	if err != nil {
		logger.Error("error parsing installation id", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Error parsing installation id: %v", err)
		return
	}

	action, _ := payload["action"].(string)
	logger.Debug("handling webhook",
		slog.Any("action", action),
		slog.Int64("installation_id", installationID),
	)

	// Check type of webhook event
	event := r.Header.Get("X-GitHub-Event")
	switch event {
	case "installation":
		installation := payload["installation"].(map[string]interface{})
		login := installation["account"].(map[string]interface{})["login"].(string)

		// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=created#installation
		if payload["action"] == "created" {
			repositories := payload["repositories"].([]interface{})

			account := installation["account"].(map[string]interface{})
			userID, _ := account["id"].(json.Number).Int64()

			logger.Info("app installation created", slog.Any("id", installation["id"]), slog.String("login", login), slog.Int("repositories", len(repositories)))

			repos := mapRepos(userID, repositories)
			err = h.repoRepo.Create(r.Context(), repos)
			if err != nil {
				logger.Error("error creating repositories", slog.Any("error", err))
				w.WriteHeader(http.StatusInternalServerError)
				break
			}
		}
		if payload["action"] == "deleted" {
			logger.Info("app installation deleted", slog.Any("id", installation["id"]), slog.String("login", login))

			// TODO: Delete all repos for this user
		}
	case "installation_repositories":
		userID, _ := payload["sender"].(map[string]interface{})["id"].(json.Number).Int64()

		// https://docs.github.com/en/webhooks/webhook-events-and-payloads#installation_repositories
		if payload["action"] == "added" {
			addedRepositories := payload["repositories_added"].([]interface{})
			repos := mapRepos(userID, addedRepositories)
			err = h.repoRepo.Create(r.Context(), repos)
			if err != nil {
				logger.Error("error creating repositories", slog.Any("error", err))
			}
		}

		if payload["action"] == "removed" {
			removedRepositories := payload["repositories_removed"].([]interface{})
			repos := mapRepos(userID, removedRepositories)
			repoIDs := make([]int64, 0, len(repos))
			for _, repo := range repos {
				repoIDs = append(repoIDs, repo.ID)
			}

			err = h.repoRepo.Delete(r.Context(), repoIDs)
			if err != nil {
				logger.Error("error deleting repositories", slog.Any("error", err))
			}
		}
	case "check_suite":
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads#check_suite
		if payload["action"] == "requested" || payload["action"] == "rerequested" {
			repository := payload["repository"].(map[string]interface{})
			repoName := repository["name"].(string)
			repoOwner := repository["owner"].(map[string]interface{})["login"].(string)
			repoID, _ := repository["id"].(json.Number).Int64()

			checkSuite := payload["check_suite"].(map[string]interface{})

			headCommit := checkSuite["head_commit"].(map[string]interface{})
			headSHA := checkSuite["head_sha"].(string)
			message := headCommit["message"].(string)

			logger.Debug("check suite requested", slog.String("owner", repoOwner), slog.String("repo", repoName), slog.String("head_sha", headSHA))

			// Create 3 random builds
			h.worker.Add(data.NewBuild{
				RepoID:    repoID,
				CommitSHA: headSHA,
				CommitMsg: message,
			})

			// Create 3 random builds
			h.worker.Add(data.NewBuild{
				RepoID:    repoID,
				CommitSHA: headSHA,
				CommitMsg: message,
			})

			// Create 3 random builds
			h.worker.Add(data.NewBuild{
				RepoID:    repoID,
				CommitSHA: headSHA,
				CommitMsg: message,
			})
		}

	default:
		logger.Error("unknown event", slog.String("event", event))
	}
}

func mapRepos(userID int64, repositories []interface{}) []data.Repo {
	repos := make([]data.Repo, 0, len(repositories))
	for _, repo := range repositories {
		repoID, _ := repo.(map[string]interface{})["id"].(json.Number).Int64()
		repos = append(repos, data.Repo{
			ID:     repoID,
			Name:   repo.(map[string]interface{})["name"].(string),
			UserID: userID,
		})
	}
	return repos
}
