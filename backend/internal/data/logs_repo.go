package data

import "context"

type LogsRepo interface {
	// Get returns all logs for the build with buildID.
	Get(ctx context.Context, buildID int64) (logs []string, err error)
}

type InfluxLogsRepo struct{}

func (r InfluxLogsRepo) Get(ctx context.Context, buildID int64) (logs []string, err error) {
	logs = []string{
		"2021-09-01T00:00:00Z: build started",
		"2021-09-01T00:00:01Z: running tests",
		"2021-09-01T00:00:02Z: tests passed",
	}

	return logs, nil
}

var _ LogsRepo = InfluxLogsRepo{}

func NewInfluxLogsRepo() *InfluxLogsRepo {
	return &InfluxLogsRepo{}
}
