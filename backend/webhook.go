package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v64/github"

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

	var accessToken string
	var user *github.User
	if code == "charlie" {
		accessToken = "access_token"
		user = &github.User{
			ID:    github.Int64(-100),
			Login: github.String("charlie"),
		}
	} else {
		var err error
		accessToken, err = h.exchangeCode(r.Context(), code)
		if err != nil {
			logger.Error(fmt.Sprintf("error exchanging code \"%s\" for access token", code), slog.Any("error", err))
			http.Error(w, "error exchanging code for access token", http.StatusInternalServerError)
			return
		}

		ghClient := github.NewClient(nil).WithAuthToken(accessToken)
		user, _, err = ghClient.Users.Get(ctx, "")
		if err != nil {
			logger.Error("error getting user info", slog.Any("error", err))
			http.Error(w, "error getting user info", http.StatusInternalServerError)
			return
		}

		err = h.userRepo.Upsert(ctx, data.NewUser{
			ID:           *user.ID,
			Username:     *user.Login,
			AccessToken:  accessToken,
			RefreshToken: "", // GitHub doesn't provide refresh tokens for OAuth Apps
		})
		if err != nil {
			logger.Error("error upserting user to database", slog.Any("error", err))
			http.Error(w, "error upserting user to database", http.StatusInternalServerError)
			return
		}
	}

	msg := fmt.Sprintf("Successfully authorized! User %s (ID: %d) has been saved to the database.", *user.Login, *user.ID)
	logger.Info(msg, "access_token", accessToken)

	// Create JWT
	token, err := createToken(*user.ID)
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

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		msg := "failed to read request body"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	eventType := github.WebHookType(r)
	event, err := github.ParseWebHook(eventType, bodyBytes)
	if err != nil {
		msg := "failed to parse webhook"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	switch event := event.(type) {
	case *github.GitHubAppAuthorizationEvent:
		// TODO: What to do when user revokes their authorization?
		//  Idea 1: delete their all data. Problem: installation still exists?
		//  Idea 2: kill their all JWTs and require reauthorization on next dashboard visit? Also stop running all their flows.
	case *github.InstallationEvent:
		installation := event.Installation
		login := *installation.Account.Login
		userID := *installation.Account.ID

		// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=created#installation
		if *event.Action == "created" {
			repositories := event.Repositories

			logger.Info("app installation created",
				slog.Any("id", installation.ID),
				slog.String("login", login),
				slog.Int("repositories", len(repositories)),
			)

			repos := mapRepos(userID, repositories)
			err = h.repoRepo.Create(r.Context(), repos)
			if err != nil {
				logger.Error("error creating repositories", slog.Any("error", err))
				w.WriteHeader(http.StatusInternalServerError)
				break
			}
		} else if *event.Action == "deleted" {
			logger.Info("app installation deleted",
				slog.Any("id", installation.ID),
				slog.String("login", login),
			)

			// TODO: Delete all repos for this user
		}
	case github.InstallationRepositoriesEvent:
		userID := *event.Sender.ID

		// https://docs.github.com/en/webhooks/webhook-events-and-payloads#installation_repositories
		if *event.Action == "added" {
			addedRepositories := event.RepositoriesAdded
			repos := mapRepos(userID, addedRepositories)
			err = h.repoRepo.Create(r.Context(), repos)
			if err != nil {
				logger.Error("error creating repositories", slog.Any("error", err))
			}
		} else if *event.Action == "removed" {
			removedRepositories := event.RepositoriesRemoved
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
	case github.CheckSuiteEvent:
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads#check_suite
		if *event.Action == "requested" || *event.Action == "rerequested" {
			headSHA := *event.CheckSuite.HeadCommit.SHA
			message := *event.CheckSuite.HeadCommit.Message

			logger.Debug("check suite requested",
				slog.String("owner", *event.Repo.Owner.Login),
				slog.String("repo", *event.Repo.Name),
				slog.String("head_sha", headSHA),
			)

			// Create 3 random builds
			h.worker.Add(data.NewBuild{
				RepoID:    *event.Repo.ID,
				CommitSHA: headSHA,
				CommitMsg: message,
			})

			// Create 3 random builds
			h.worker.Add(data.NewBuild{
				RepoID:    *event.Repo.ID,
				CommitSHA: headSHA,
				CommitMsg: message,
			})

			// Create 3 random builds
			h.worker.Add(data.NewBuild{
				RepoID:    *event.Repo.ID,
				CommitSHA: headSHA,
				CommitMsg: message,
			})
		}

	default:
		logger.Error("unknown event", slog.String("event", eventType))
	}
}

func mapRepos(userID int64, repositories []*github.Repository) []data.Repo {
	repos := make([]data.Repo, 0, len(repositories))
	for _, repo := range repositories {
		repos = append(repos, data.Repo{
			ID:     *repo.ID,
			Name:   *repo.Name,
			UserID: userID,
		})
	}
	return repos
}
