package listener

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/bartekpacia/ghapp/data"

	"github.com/lib/pq"
)

const channelName = "builds_channel"

type Listener struct {
	db          *sql.DB
	listener    *pq.Listener
	channelName string
	logger      *slog.Logger
}

func NewListener(db *sql.DB, connInfo string) *Listener {
	minReconn := 10 * time.Second
	maxReconn := time.Minute
	listener := pq.NewListener(connInfo, minReconn, maxReconn, nil)

	return &Listener{
		db:          db,
		channelName: channelName,
		listener:    listener,
		logger:      slog.Default(), // TODO: add some "subsystem name" to this logger
	}
}

func (l Listener) Start(ctx context.Context) error {
	err := l.listener.Listen(channelName)
	if err != nil {
		return fmt.Errorf("listen on channel %s: %w", channelName, err)
	}

	l.logger.Info("started listener", slog.String("channel", channelName))

	for {
		select {
		case <-ctx.Done():
			l.logger.Debug("context cancelled, stopping listener")
			_ = l.listener.Close()
			return nil
		case msg := <-l.listener.Notify:
			l.logger.Debug("received notification", slog.Any("channel", msg.Channel))

			updatedBuild := data.Build{}
			err := json.Unmarshal([]byte(msg.Extra), &updatedBuild)
			if err != nil {
				// TODO: handle error
				l.logger.Error("failed to unmarshal build", slog.Any("error", err))
			}
		}
	}
}

// Returns immediately and starts a goroutine in the background
func (l Listener) createCheckRun(ctx context.Context, build data.Build) error {
	logger := l.logger

	// FIXME: FIMXE!!
	// TODO: go from thread to ball...
	owner := "bartekpacia"
	repo := "dumbpkg"

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/check-runs", owner, repo)

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
