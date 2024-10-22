package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	buildID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	influxURL := mustGetenv("INFLUXDB_URL")
	influxToken := mustGetenv("INFLUXDB_TOKEN")
	influxBucket := mustGetenv("INFLUXDB_BUCKET")
	influxOrg := mustGetenv("INFLUXDB_ORG")
	influxClient := influxdb2.NewClient(influxURL, influxToken)
	_, err = influxClient.Health(ctx)
	if err != nil {
		slog.Error("error connecting to Influx database", slog.Any("error", err))
		os.Exit(1)

	}
	slog.Info("connected to Influx database", "url", influxURL)

	query := fmt.Sprintf("from(bucket: \"%s\") |> range(start: -1h) |> filter(fn: (r) => r[\"_measurement\"] == \"%d\")", influxBucket, buildID)
	fmt.Println("will execute query:", query)
	queryAPI := influxClient.QueryAPI(influxOrg)
	queryResult, err := queryAPI.Query(ctx, query)
	if err != nil {
		slog.Error("query influxdb", slog.Any("error", err))
		os.Exit(1)
	}

	logs := make([]string, 0)
	for queryResult.Next() {
		record := queryResult.Record()
		logs = append(logs, fmt.Sprintf("%s: %s", record.Time(), record.Value()))
	}

	slog.Info("success! retrieved logs")

	for _, log := range logs {
		fmt.Println(log)
	}
}

func mustGetenv(varname string) string {
	value := os.Getenv(varname)
	if value == "" {
		slog.Error(varname + " env var is empty or not set")
		os.Exit(1)
	}
	return value
}
