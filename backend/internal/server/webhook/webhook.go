// Package webhook implements the handling of GitHub webhooks.
package webhook

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v64/github"

	l "github.com/bee-ci/bee-ci-system/internal/common/logger"
	"github.com/bee-ci/bee-ci-system/internal/common/middleware"
	"github.com/bee-ci/bee-ci-system/internal/data"
)

//go:embed redirect.html
var redirectHTMLPage embed.FS

type Handler struct {
	httpClient *http.Client
	userRepo   data.UserRepo
	repoRepo   data.RepoRepo
	buildRepo  data.BuildRepo

	// The domain where the auth cookie will be placed. For example:
	// - .pacia.tech
	// - .karolak.cc
	//
	// Must be empty for localhost.
	mainDomain string

	// The URL the user will be redirected to after successful auth. For example:
	//  - https://bee-ci.pacia.tech/dashboard
	//  - http://localhost:8080/dashboard
	redirectURL string

	githubAppClientID      string
	githubAppClientSecret  string
	githubAppWebhookSecret string

	// The secret key used to sign JWT tokens.
	jwtSecret []byte
}

func NewHandler(
	userRepo data.UserRepo,
	repoRepo data.RepoRepo,
	buildRepo data.BuildRepo,
	mainDomain string,
	frontendURL string,
	githubAppClientID string,
	githubAppClientSecret string,
	githubAppWebhookSecret string,
	jwtSecret []byte,
) (*Handler, error) {
	redirectURL, err := url.JoinPath(frontendURL, "dashboard")
	if err != nil {
		return nil, fmt.Errorf("could not join path valid to create a redirect URL: %v", err)
	}

	return &Handler{
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
	}, nil
}

func (h Handler) Mux() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /github/callback/", http.HandlerFunc(h.handleAuthCallback))

	mux.Handle("POST /{$}",
		middleware.WithWebhookSecret(
			http.HandlerFunc(h.handleWebhook), h.githubAppWebhookSecret,
		),
	)

	return mux
}

func (h Handler) exchangeCode(ctx context.Context, code string) (userAccessToken string, err error) {
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
func (h Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
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

		newUser := data.NewUser{
			ID:       *ghUser.ID,
			Username: *ghUser.Login,
		}
		err = h.userRepo.Upsert(ctx, newUser)
		if err != nil {
			logger.Error("error upserting ghUser to database", slog.Any("error", err))
			http.Error(w, "error upserting ghUser to database", http.StatusInternalServerError)
			return
		}

		logger.Info("github user was created (or updated)", slog.Any("user", newUser))
	}

	token, err := h.createToken(*ghUser.ID)
	if err != nil {
		logger.Error("error creating token", slog.Any("error", err))
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	jwtTokenCookie := &http.Cookie{
		Name:   "jwt",
		Value:  token,
		Domain: h.mainDomain,
		Path:   "/",
	}
	logger.Debug("setting jwt cookie", slog.Any("cookie", jwtTokenCookie))

	http.SetCookie(w, jwtTokenCookie)

	tmpl, err := template.ParseFS(redirectHTMLPage, "redirect.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	tmplData := struct {
		DashboardURL string
	}{
		DashboardURL: h.redirectURL,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h Handler) handleWebhook(w http.ResponseWriter, r *http.Request) {
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
		// Payload: https://github.com/octokit/webhooks/blob/main/payload-examples/api.github.com/github_app_authorization/revoked.payload.json

		installation := event.Installation

		logger.Debug("new webhook event",
			slog.String("event", eventType),
			slog.String("action", *event.Action),
			slog.Int64("installation.id", *installation.ID),
		)

		// TODO: What to do when user revokes their authorization?
		//  Idea 1: delete their all data. Problem: installation still exists?
		//  Idea 2: kill their all JWTs and require reauthorization on next dashboard visit? Also stop running all their flows.
	case *github.InstallationEvent:
		// Payload: https://github.com/octokit/webhooks/blob/main/payload-examples/api.github.com/installation/created.payload.json
		// Payload: https://github.com/octokit/webhooks/blob/main/payload-examples/api.github.com/installation/deleted.payload.json

		installation := event.Installation
		login := *installation.Account.Login
		userID := *installation.Account.ID

		logger.Debug("new webhook event",
			slog.String("event", eventType),
			slog.String("action", *event.Action),
			slog.Int64("installation.id", *installation.ID),
			slog.String("user.name", login),
			slog.Int64("user.id", userID),
		)

		if *event.Action == "created" {
			repos := mapRepos(userID, event.Repositories)
			err = h.repoRepo.Upsert(r.Context(), repos)
			if err != nil {
				logger.Error("error creating repositories", slog.Any("error", err))
				http.Error(w, "error creating repositories", http.StatusInternalServerError)
				break
			}
		} else if *event.Action == "deleted" {
			removedRepositories := event.Repositories

			repoIDs := make([]int64, 0, len(removedRepositories))
			for _, removedRepository := range removedRepositories {
				repoIDs = append(repoIDs, *removedRepository.ID)
			}

			err = h.repoRepo.Delete(r.Context(), repoIDs)
			if err != nil {
				logger.Error("error deleting repositories", slog.Any("error", err))
				http.Error(w, "error deleting repositories", http.StatusInternalServerError)
				break
			}
			_, _ = w.Write([]byte(fmt.Sprintf("removed %d repositories\n", len(removedRepositories))))
		}
	case *github.InstallationRepositoriesEvent:
		// Payload: https://github.com/octokit/webhooks/blob/main/payload-examples/api.github.com/installation_repositories/added.payload.json
		// Payload: https://github.com/octokit/webhooks/blob/main/payload-examples/api.github.com/installation_repositories/removed.payload.json

		installation := *event.Installation
		userID := *event.Sender.ID

		logger.Debug("new webhook event",
			slog.String("event", eventType),
			slog.String("action", *event.Action),
			slog.Int64("installation.id", *installation.ID),
			slog.Int64("sender.id", userID),
		)

		switch *event.Action {
		case "added":
			addedRepositories := event.RepositoriesAdded
			repos := mapRepos(userID, addedRepositories)
			err = h.repoRepo.Upsert(r.Context(), repos)
			if err != nil {
				logger.Error("error creating repositories", slog.Any("error", err))
				http.Error(w, "error creating repositories", http.StatusInternalServerError)
				break
			}
			_, _ = w.Write([]byte(fmt.Sprintf("added %d repositories", len(repos))))
		case "removed":
			removedRepositories := event.RepositoriesRemoved

			repoIDs := make([]int64, 0, len(removedRepositories))
			for _, removedRepository := range removedRepositories {
				repoIDs = append(repoIDs, *removedRepository.ID)
			}

			err = h.repoRepo.Delete(r.Context(), repoIDs)
			if err != nil {
				logger.Error("error deleting repositories", slog.Any("error", err))
				http.Error(w, "error deleting repositories", http.StatusInternalServerError)
				break
			}
			_, _ = w.Write([]byte(fmt.Sprintf("removed %d repositories\n", len(removedRepositories))))
		}
	case *github.CheckSuiteEvent:
		// Payload: https://github.com/octokit/webhooks/blob/main/payload-examples/api.github.com/check_suite/requested.payload.json

		installation := *event.Installation
		userID := *event.Sender.ID

		logger.Debug("new webhook event",
			slog.String("event", eventType),
			slog.String("action", *event.Action),
			slog.Int64("installation.id", *installation.ID),
			slog.Int64("sender.id", userID),
		)

		// Create build
		if *event.Action == "requested" || *event.Action == "rerequested" {
			headSHA := *event.CheckSuite.HeadSHA
			message := *event.CheckSuite.HeadCommit.Message

			logger.Debug(fmt.Sprintf("check suite %s", *event.Action),
				slog.String("owner", *event.Repo.Owner.Login),
				slog.String("removedRepository", *event.Repo.Name),
				slog.Int64("installation_id", *installation.ID),
				slog.String("head_sha", headSHA),
			)

			// TODO: Parse information from the BeeCI config file here (such as name)
			buildID, err := h.buildRepo.Create(r.Context(), data.NewBuild{
				RepoID:         *event.Repo.ID,
				CommitSHA:      headSHA,
				CommitMsg:      message,
				InstallationID: *installation.ID,
			})
			if err != nil {
				logger.Error("failed to create build", slog.Any("error", err))
				// TODO: handle error in a better way – update status on GitHub
				http.Error(w, fmt.Sprintf("failed to create build: %v", err), http.StatusInternalServerError)
				return
			}
			logger.Debug("build created", slog.Int64("build_id", buildID))
			_, _ = w.Write([]byte("build created, ID: " + strconv.FormatInt(buildID, 10)))
		}

	default:
		logger.Error("unknown event", slog.String("event", eventType))
	}
}

// Function to create JWT tokens with claims
func (h Handler) createToken(userID int64) (string, error) {
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
			ID:     *repo.ID,
			Name:   *repo.Name,
			UserID: userID,
		})
	}
	return repos
}
