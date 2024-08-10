package queue

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bartekpacia/ghapp/data"
	"log/slog"
	"time"

	"github.com/lib/pq"
)

const channelName = "builds_channel"

type Listener struct {
	db            *sql.DB
	channelName   string
	notifications <-chan *pq.Notification
	l             *slog.Logger
}

func NewListener(db *sql.DB, connInfo string) *Listener {
	minReconn := 10 * time.Second
	maxReconn := time.Minute

	callback := func(ev pq.ListenerEventType, err error) {
		if err != nil {

		}
	}

	listener := pq.NewListener(connInfo, minReconn, maxReconn, callback)
	err := listener.Listen(channelName)
	if err != nil {
		panic(err)
	}

	return &Listener{
		db:            db,
		channelName:   channelName,
		notifications: listener.NotificationChannel(),
		l:             slog.Default(), // TODO: add some "subsystem name" to this logger
	}
}

func (l Listener) Start() {
	for {
		select {
		case msg := <-l.notifications:
			updatedBuild := data.Build{}
			// decoder := json.NewDecoder(bytes.NewBufferString(msg.Extra))

			err := json.Unmarshal([]byte(msg.Extra), &updatedBuild)
			if err != nil {
				// handle error
				l.l.Error("failed to unmarshal build", slog.Any("error", err))
			}
			l.l.Info("received notification", slog.Any("build", updatedBuild))
			fmt.Printf("same but with more info: %+v\n", updatedBuild)
		}
	}
}
