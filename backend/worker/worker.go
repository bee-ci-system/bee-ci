// Package worker implements a worker that executes jobs.
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

type Worker interface {
	Add(build data.NewBuild)
}

type worker struct {
	ctx       context.Context
	buildRepo data.BuildRepo
}

func (w worker) Add(build data.NewBuild) {
	go w.job(build)
}

func (w worker) job(build data.NewBuild) {
	slog.Info("Starting job for build")
	// Do some work
	buildId, err := w.buildRepo.Create(w.ctx, build)
	if err != nil {
		slog.Error("failed to create build", slog.Any("error", err))
		// TODO: handle error
		return
	}

	slog.Info("job queued", slog.Uint64("build_id", buildId))

	time.Sleep(5 * time.Second)
	err = w.buildRepo.Update(w.ctx, buildId, data.StatusRunning)
	if err != nil {
		slog.Error("failed to update build status", slog.Any("error", err))
		return
	}

	slog.Info("job running", slog.Uint64("build_id", buildId))

	time.Sleep(5 * time.Second)

	// random failure or success, 50% chance of failure
	if rand.Intn(2) == 0 {
		w.buildRepo.Update(w.ctx, buildId, data.StatusFailed)
		slog.Info("job failed", slog.Uint64("build_id", buildId))
	} else {
		w.buildRepo.Update(w.ctx, buildId, data.StatusSuccess)
		slog.Info("job succeeded", slog.Uint64("build_id", buildId))
	}
}

func New(ctx context.Context, buildRepo data.BuildRepo) Worker {
	return &worker{
		ctx:       ctx,
		buildRepo: buildRepo,
	}
}
