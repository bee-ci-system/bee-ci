package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/bartekpacia/ghapp/data"
	"github.com/bartekpacia/ghapp/worker"
)

type WebhookHandler struct {
	userRepo   data.UserRepo
	httpClient *http.Client
	worker     worker.Worker
}

func NewWebhookHandler(userRepo data.UserRepo, w worker.Worker) *WebhookHandler {
	return &WebhookHandler{
		userRepo: userRepo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		worker: w,
	}
}

func (h WebhookHandler) Mux() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", http.HandlerFunc(handleIndex))
	mux.Handle("GET /github/callback", http.HandlerFunc(h.handleAuthCallback))

	mux.HandleFunc("POST /{$}", func (w http.ResponseWriter, r *http.Request) {
		log.Println("HEHEHE HERE!")
	})

	mux.Handle("POST /",
		WithWebhookSecret(
			WithAuthenticatedApp( // provides gh_app_client
				WithAuthenticatedAppInstallation( // provides gh_installation_client
					http.HandlerFunc(h.handleWebhook),
				),
			),
		),
	)

	return mux
}

func (h WebhookHandler) githubUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	l := ctx.Value(ctxLogger{}).(*slog.Logger)
	const url = "https://api.github.com/user"

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
	l.Info("authenticated user",
		slog.Any("login", userData["login"]),
		slog.Any("id", userData["id"]),
		slog.Any("name", userData["name"]),
	)

	return userData, nil
}

func (h WebhookHandler) exchangeCode(ctx context.Context, code string) (userAccessToken string, err error) {
	const url = "https://github.com/login/oauth/access_token"

	reqBody := map[string]interface{}{
		"client_id":     clientID,
		"client_secret": clientSecret,
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
	l := r.Context().Value(ctxLogger{}).(*slog.Logger)

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code query parameter", http.StatusBadRequest)
		return
	}

	accessToken, err := h.exchangeCode(r.Context(), code)
	if err != nil {
		l.Error("error exchanging code for access token", slog.Any("error", err))
		http.Error(w, "error exchanging code for access token", http.StatusInternalServerError)
		return
	}

	_, err = h.githubUserInfo(r.Context(), accessToken)
	if err != nil {
		l.Error("error getting user info", slog.Any("error", err))
		http.Error(w, "error getting user info", http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Successfully authorized! Got code %s and exchanged it for a user access token ending in %s", code, accessToken[len(accessToken)-9:])
	l.Info(msg)

	_, _ = fmt.Fprint(w, msg)
}

func (h WebhookHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	l := r.Context().Value(ctxLogger{}).(*slog.Logger)

	event := r.Header.Get("X-GitHub-Event")

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	var payload map[string]interface{}
	err := decoder.Decode(&payload)
	if err != nil {
		msg := "error reading request body"
		l.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	installationIDStr := payload["installation"].(map[string]interface{})["id"].(json.Number).String()
	installationID, err := strconv.ParseInt(installationIDStr, 10, 64)
	if err != nil {
		l.Error("error parsing installation id", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Error parsing installation id: %v", err)
		return
	}

	action, _ := payload["action"].(string)
	l.Info("handling webhook",
		slog.String("event", event),
		slog.Any("action", action),
		slog.Int64("installation_id", installationID),
	)

	// Check type of webhook event
	switch event {
	case "installation":
		installation := payload["installation"].(map[string]interface{})
		login := installation["account"].(map[string]interface{})["login"].(string)

		// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=created#installation
		if payload["action"] == "created" {
			repositories := payload["repositories"].([]interface{})
			l.Info("app installation created", slog.Any("id", installation["id"]), slog.String("login", login), slog.Int("repositories", len(repositories)))
		}
		if payload["action"] == "deleted" {
			l.Info("app installation deleted", slog.Any("id", installation["id"]), slog.String("login", login))
		}
	case "check_suite":
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=requested#check_suite
		if payload["action"] == "requested" || payload["action"] == "rerequested" {
			repository := payload["repository"].(map[string]interface{})
			repoName := repository["name"].(string)
			repoOwner := repository["owner"].(map[string]interface{})["login"].(string)

			checkSuite := payload["check_suite"].(map[string]interface{})
			headSHA := checkSuite["head_sha"].(string)
			l.Info("check suite requested", slog.String("owner", repoOwner), slog.String("repo", repoName), slog.String("head_sha", headSHA))

			// Create 3 builds: both run for 10 seconds, then 1 fails, 1 succeeds, 1 never stops (always "queued")
			go func() {
				err := h.createCheckRun(r.Context(), repoOwner, repoName, headSHA, "i will pass", "success")
				if err != nil {
					l.Error("error creating check run", slog.Any("error", err))
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprintf(w, "Error creating check run: %v", err)
				}
			}()

			go func() {
				err = h.createCheckRun(r.Context(), repoOwner, repoName, headSHA, "i will fail", "failure")
				if err != nil {
					l.Error("error creating check run", slog.Any("error", err))
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprintf(w, "Error creating check run: %v", err)
				}
			}()

			go func() {
				err = h.createCheckRun(r.Context(), repoOwner, repoName, headSHA, "i will keep running", "queued")
				if err != nil {
					l.Error("error creating check run", slog.Any("error", err))
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprintf(w, "Error creating check run: %v", err)
				}
			}()
		}
	//case "check_run":
	//	// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=created#check_run
	//	repository := payload["repository"].(map[string]interface{})
	//	repoName := repository["name"].(string)
	//	repoOwner := repository["owner"].(map[string]interface{})["login"].(string)
	//
	//	checkRun := payload["check_run"].(map[string]interface{})
	//	headSHA := checkRun["head_sha"].(string)
	//	err := createCheckRun(r.Context(), repoOwner, repoName, headSHA)
	//	if err != nil {
	//		msg := "error creating check run"
	//		l.Error(msg, slog.Any("error", err))
	//		http.Error(w, msg, http.StatusInternalServerError)
	//		return
	//	}
	default:
		l.Error("unknown event", slog.String("event", event))
	}
}

// TODO: accept context, and access logger and authenticated HTTP client from there?

// Returns immediately and starts a goroutine in the background
func (h WebhookHandler) createCheckRun(ctx context.Context, owner, repo, sha string, msg string, conclusion string) error {
	// START: NEW PORT THAT WRITES TO WORKER

	h.worker.Add(data.NewBuild{
		RepoID:    1,
		CommitSHA: sha,
	})

	// END: NEW PORT THAT WRITES TO WORKER

	l := ctx.Value(ctxLogger{}).(*slog.Logger)
	githubInstallationClient := ctx.Value(ctxGHInstallationClient{}).(http.Client)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/check-runs", owner, repo)

	body := map[string]interface{}{
		"head_sha":    sha,
		"name":        msg + ", started at: " + fmt.Sprint(time.Now().Format(time.RFC822Z)),
		"details_url": "https://garden.pacia.com",
		"status":      "in_progress",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling body to JSON: %w", err)
	}

	// Don't use context ctx here, because it'll get canceled by caller
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	res, err := githubInstallationClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending POST request: %w", err)
	}

	respBody := make([]byte, 0)
	_, err = res.Body.Read(respBody)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	l.Info("initial request made", slog.Int("status", res.StatusCode), slog.String("body", string(respBody)))

	time.Sleep(10 * time.Second)
	switch conclusion {
	case "success":
		body["status"] = "completed"
		body["conclusion"] = "success"
	case "failure":
		body["status"] = "completed"
		body["conclusion"] = "failure"
	default:
		// no conclusion, keep running this
	}

	bodyBytes, err = json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling body to JSON: %w", err)
	}

	req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	res, err = githubInstallationClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending POST request: %w", err)
	}

	l.Info("final request made", slog.Int("status", res.StatusCode), slog.String("body", string(respBody)))

	return nil
}
