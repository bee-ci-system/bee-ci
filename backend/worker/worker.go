// Package worker implements a Worker that executes jobs.
//
// It spawns a single goroutine per new job.
package worker

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/bartekpacia/ghapp/data"
)

type Worker struct {
	ctx       context.Context
	buildRepo data.BuildRepo
}

func New(ctx context.Context, buildRepo data.BuildRepo) *Worker {
	return &Worker{
		ctx:       ctx,
		buildRepo: buildRepo,
	}
}

func (w Worker) Add(build data.NewBuild) {
	go w.job(build)
}

func (w Worker) job(build data.NewBuild) {
	slog.Info("Starting job for build")
	buildId, err := w.buildRepo.Create(w.ctx, build)
	if err != nil {
		slog.Error("failed to create build", slog.Any("error", err))
		// TODO: handle error
		return
	}

	slog.Info("job queued", slog.Int64("build_id", buildId))

	time.Sleep(5 * time.Second)
	err = w.buildRepo.UpdateStatus(w.ctx, buildId, "in_progress")
	if err != nil {
		slog.Error("failed to update build status", slog.Any("error", err))
		return
	}

	slog.Info("job running", slog.Int64("build_id", buildId))

	time.Sleep(5 * time.Second)

	// random failure or success, 50% chance of failure
	if rand.Intn(2) == 0 {
		err := w.buildRepo.SetConclusion(w.ctx, buildId, "failure")
		if err != nil {
			slog.Error("failed to set failure conclusion", slog.Any("error", err))
			return
		}

		slog.Info("job failed", slog.Int64("build_id", buildId))
	} else {
		err := w.buildRepo.SetConclusion(w.ctx, buildId, "success")
		if err != nil {
			slog.Error("failed to update success conclusion", slog.Any("error", err))
			return
		}
		slog.Info("job succeeded", slog.Int64("build_id", buildId))
	}
}
