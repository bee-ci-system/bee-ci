// Package updater implements a listener that listens the database for build
// updates and creates/updates check runs on GitHub.
package updater

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/bee-ci/bee-ci-system/data"

	"github.com/lib/pq"
)

const channelName = "builds_channel"

type Updater struct {
	listener    *pq.Listener
	channelName string
	logger      *slog.Logger
	repoRepo    data.RepoRepo
	userRepo    data.UserRepo
	buildRepo   data.BuildRepo
}

func New(dbListener *pq.Listener, repoRepo data.RepoRepo, userRepo data.UserRepo, buildRepo data.BuildRepo) *Updater {
	listener := dbListener

	return &Updater{
		channelName: channelName,
		listener:    listener,
		logger:      slog.Default(), // TODO: add some "subsystem name" to this logger
		repoRepo:    repoRepo,
		userRepo:    userRepo,
		buildRepo:   buildRepo,
	}
}

// Start starts the updater. It will listen for updates from the database and
// create check runs on GitHub when the updates happen.
//
// To shutdown the updater, cancel the context.
func (l Updater) Start(ctx context.Context) error {
	err := l.listener.Listen(channelName)
	if err != nil {
		return fmt.Errorf("listen on channel %s: %w", channelName, err)
	}

	l.logger.Info("started listener", slog.String("channel", channelName))

	for {
		select {
		case <-ctx.Done():
			l.logger.Debug("context cancelled, stopping listener")
			err = l.listener.Close()
			if err != nil {
				l.logger.Error("failed to close listener", slog.Any("error", err))
				return err
			}
			return nil
		case msg := <-l.listener.Notify:
			l.logger.Debug("received notification", slog.Any("channel", msg.Channel))

			updatedBuild := data.Build{}
			err := json.Unmarshal([]byte(msg.Extra), &updatedBuild)
			if err != nil {
				l.logger.Error("failed to unmarshal build", slog.Any("error", err))
				break
			}

			err = l.createCheckRun(ctx, updatedBuild)
			if err != nil {
				l.logger.Error("failed to create check run", slog.Any("error", err))
				break
			}

			l.logger.Info("check run created", slog.Any("build", updatedBuild))
		}
	}
}

// Returns immediately and starts a goroutine in the background
func (l Updater) createCheckRun(ctx context.Context, build data.Build) error {
	logger := l.logger

	repo, err := l.repoRepo.Get(ctx, build.RepoID)
	if err != nil {
		return fmt.Errorf("getting repo: %w", err)
	}
	repoName := repo.Name

	user, err := l.userRepo.Get(ctx, repo.UserID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}
	owner := user.Username

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/check-runs", owner, repoName)

	// githubInstallationClient := ctx.Value(ctxGHInstallationClient{}).(http.Client)
	githubInstallationClient := http.DefaultClient

	body := map[string]interface{}{
		"external_id": build.ID,
		"head_sha":    build.CommitSHA,
		"name":        "build.CommitMSG" + ", started at: " + fmt.Sprint(time.Now().Format(time.RFC822Z)),
		"details_url": "https://garden.pacia.com",
		"status":      build.Status,
	}
	if build.Conclusion != nil {
		body["conclusion"] = build.Conclusion
		body["completed_at"] = build.UpdatedAt
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling body to JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := githubInstallationClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBodyBytes := make([]byte, 0)
	_, err = resp.Body.Read(respBodyBytes)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	logger.Info("request made", slog.Int("status", resp.StatusCode), slog.String("body", string(respBodyBytes)))

	return nil
}
