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

var l = slog.Default() // TODO: add some "subsystem name" to this logger

type Queue struct {
	db *sql.DB
}

type Listener struct {
	db            *sql.DB
	channelName   string
	notifications <-chan pq.Notification
}

func NewListener(connInfo string) *Listener {
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

}

func (listener Listener) Start() {
	for {
		select {
		case msg := <-listener.notifications:
			updatedBuild := data.Build{}
			err := json.Unmarshal([]byte(msg.Extra), &updatedBuild)
			if err != nil {
				// handle error
				l.Error("failed to unmarshal build", slog.Any("error", err))
			}
			l.Info("received notification", slog.Any("build", updatedBuild))
			fmt.Printf("same but with more info: %+v\n", updatedBuild)
		}
	}
}

func New(db *sql.DB) *Queue {
	return &Queue{
		db: db,
		// Listener: newListener(db, channelName),
	}
}

// Start starts the queue:
//  1. Listens for changes in the builds table in the database
//  2. Sends the build status to GitHub
func (q Queue) Start() {
	pq.NewListener()
}
