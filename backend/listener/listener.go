package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
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

	for {
		select {
		case <-ctx.Done():
			l.logger.Info("context cancelled, stopping listener")
			_ = l.listener.Close()
			return nil
		case msg := <-l.listener.Notify:
			l.logger.Info("received notification", slog.Any("channel", msg.Channel))

			updatedBuild := data.Build{}
			err := json.Unmarshal([]byte(msg.Extra), &updatedBuild)
			if err != nil {
				// TODO: handle error
				l.logger.Error("failed to unmarshal build", slog.Any("error", err))
			}
		}
	}
}
