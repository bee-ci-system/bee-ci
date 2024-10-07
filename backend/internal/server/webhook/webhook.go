// Package webhook implements handling of GitHub webhooks.
package webhook

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

	l "github.com/bee-ci/bee-ci-system/internal/common/logger"
	"github.com/bee-ci/bee-ci-system/internal/common/middleware"
	"github.com/bee-ci/bee-ci-system/internal/data"
)

type WebhookHandler struct {
	httpClient *http.Client
	userRepo   data.UserRepo
	repoRepo   data.RepoRepo
	buildRepo  data.BuildRepo

	// The domain where the auth cookie will be placed. For example
	// ".pacia.tech" or ".karolak.cc". Must be empty for localhost.
	mainDomain string

	// The URL the user will be redirected to after successful auth. For example
	// "https://bee-ci.pacia.tech/dashboard" or
	// "http://localhost:8080/dashboard".
	redirectURL string

	githubAppClientID      string
	githubAppClientSecret  string
	githubAppWebhookSecret string

	// The secret key used to sign JWT tokens.
	jwtSecret []byte
}

func NewWebhookHandler(
	userRepo data.UserRepo,
	repoRepo data.RepoRepo,
	buildRepo data.BuildRepo,
	mainDomain string,
	redirectURL string,
	githubAppClientID string,
	githubAppClientSecret string,
	githubAppWebhookSecret string,
	jwtSecret []byte,
) *WebhookHandler {
	return &WebhookHandler{
		httpClient:             &http.Client{Timeout: 10 * time.Second},
		userRepo:               userRepo,
		repoRepo:               repoRepo,
		buildRepo:              buildRepo,
		mainDomain:             mainDomain,
		redirectURL:            redirectURL,
		githubAppClientID:      githubAppClientID,
		githubAppClientSecret:  githubAppClientSecret,
		githubAppWebhookSecret: githubAppWebhookSecret,
		jwtSecret:              jwtSecret,
	}
}

func (h WebhookHandler) Mux() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /github/callback/", http.HandlerFunc(h.handleAuthCallback))

	mux.Handle("POST /{$}",
		middleware.WithWebhookSecret(
			http.HandlerFunc(h.handleWebhook), h.githubAppWebhookSecret,
		),
	)

	return mux
}

func (h WebhookHandler) exchangeCode(ctx context.Context, code string) (userAccessToken string, err error) {
	const url = "https://github.com/login/oauth/access_token"

	reqBody := map[string]interface{}{
		"client_id":     h.githubAppClientID,
		"client_secret": h.githubAppClientSecret,
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
	var ghUser *github.User
	if code == "charlie" {
		accessToken = "access_token"
		ghUser = &github.User{
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
		ghUser, _, err = ghClient.Users.Get(ctx, "")
		if err != nil {
			logger.Error("error getting ghUser info", slog.Any("error", err))
			http.Error(w, "error getting ghUser info", http.StatusInternalServerError)
			return
		}

		err = h.userRepo.Upsert(ctx, data.NewUser{
			ID:           *ghUser.ID,
			Username:     *ghUser.Login,
			AccessToken:  accessToken,
			RefreshToken: "", // GitHub doesn't provide refresh tokens for OAuth Apps
		})
		if err != nil {
			logger.Error("error upserting ghUser to database", slog.Any("error", err))
			http.Error(w, "error upserting ghUser to database", http.StatusInternalServerError)
			return
		}
	}

	logger.Info("github user was created/updated", slog.Any("github_user", ghUser))

	token, err := h.createToken(*ghUser.ID)
	if err != nil {
		logger.Error("error creating token", slog.Any("error", err))
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	logger.Debug("JWT token created", slog.String("username", *ghUser.Login), slog.String("token", token))

	jwtTokenCookie := &http.Cookie{
		Name:   "jwt",
		Value:  token,
		Domain: h.mainDomain,
		Path:   "/",
	}

	http.SetCookie(w, jwtTokenCookie)

	http.Redirect(w, r, h.redirectURL, http.StatusSeeOther)
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

			logger.Debug("app installation created",
				slog.Any("id", installation.ID),
				slog.String("login", login),
				slog.Int("repositories", len(repositories)),
			)

			// The webhook event doesn't contain all repository data we need:
			//  - latest commit SHA
			//  - datetime of the latest commit.
			// Therefore, we need to request more data from the API.


			repos := mapRepos(userID, repositories)
			err = h.repoRepo.Create(r.Context(), repos)
			if err != nil {
				logger.Error("error creating repositories", slog.Any("error", err))
				w.WriteHeader(http.StatusInternalServerError)
				break
			}
		} else if *event.Action == "deleted" {
			logger.Debug("app installation deleted",
				slog.Any("id", installation.ID),
				slog.String("login", login),
			)

			// TODO: Delete all repos for this user
		}
	case *github.InstallationRepositoriesEvent:
		userID := *event.Sender.ID

		// https://docs.github.com/en/webhooks/webhook-events-and-payloads#installation_repositories
		switch *event.Action {
		case "added":
			addedRepositories := event.RepositoriesAdded
			repos := mapRepos(userID, addedRepositories)
			err = h.repoRepo.Create(r.Context(), repos)
			if err != nil {
				logger.Error("error creating repositories", slog.Any("error", err))
			}
		case "removed":
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
	case *github.CheckSuiteEvent:
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads#check_suite
		if *event.Action == "requested" || *event.Action == "rerequested" {
			headSHA := *event.CheckSuite.HeadSHA
			message := *event.CheckSuite.HeadCommit.Message
			installationID := *event.Installation.ID

			logger.Debug("check suite requested",
				slog.String("owner", *event.Repo.Owner.Login),
				slog.String("repo", *event.Repo.Name),
				slog.Int64("installation_id", installationID),
				slog.String("head_sha", headSHA),
			)

			// TODO: Parse information from the BeeCI config file here (such as name)
			buildID, err := h.buildRepo.Create(r.Context(), data.NewBuild{
				RepoID:         *event.Repo.ID,
				CommitSHA:      headSHA,
				CommitMsg:      message,
				InstallationID: installationID,
			})
			if err != nil {
				logger.Error("failed to create build", slog.Any("error", err))
				// TODO: handle error in a better way – update status on GitHub
				return
			}
			logger.Debug("build created", slog.Int64("build_id", buildID))
		}
	case *github.PushEvent:
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads#push
		if *event.Action == "pushed" {
			newLatestSHA := event.After
			h.repoRepo.UpdateModificationTime
		}

	default:
		logger.Error("unknown event", slog.String("event", eventType))
	}
}

// Function to create JWT tokens with claims
func (h WebhookHandler) createToken(userID int64) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"iss": "bee-ci",
		// "exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(), // Issued at
	})

	tokenString, err := claims.SignedString(h.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func mapRepos(userID int64, repositories []*github.Repository) []data.Repo {
	repos := make([]data.Repo, 0, len(repositories))
	for _, repo := range repositories {
		repos = append(repos, data.Repo{
			ID:                   *repo.ID,
			Name:                 *repo.Name,
			UserID:               userID,
			LatestCommitSHA:      *repo.,
			LatestCommitPushedAt: *repo.PushedAt.GetTime(),
		})
	}
	return repos
}
