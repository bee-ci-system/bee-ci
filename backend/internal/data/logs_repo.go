package data

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type LogsRepo interface {
	// Get returns all logs for the build with buildID.
	Get(ctx context.Context, buildID int64) (logs []string, err error)
}

type InfluxLogsRepo struct {
	influxClient influxdb2.Client
	org          string
	bucket       string
}

func (r InfluxLogsRepo) Get(ctx context.Context, buildID int64) (logs []string, err error) {
	logs = make([]string, 0)

	query := fmt.Sprintf("from(bucket: \"%s\") |> range(start: -1h) |> filter(fn: (r) => r[\"_measurement\"] == \"%d\")", r.bucket, buildID)
	queryAPI := r.influxClient.QueryAPI(r.org)
	queryResult, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query influxdb: %w", err)
	}

	for queryResult.Next() {
		record := queryResult.Record()
		logs = append(logs, fmt.Sprintf("%s: %s", record.Time(), record.Value()))
	}

	return logs, nil
}

var _ LogsRepo = InfluxLogsRepo{}

func NewInfluxLogsRepo(influxClient influxdb2.Client, org, bucket string) *InfluxLogsRepo {
	return &InfluxLogsRepo{
		influxClient: influxClient,
		org:          org,
		bucket:       bucket,
	}
}
