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
	buildId, err := w.buildRepo.Create(w.ctx, build)
	if err != nil {
		slog.Error("failed to create build", slog.Any("error", err))
		// TODO: handle error
		return
	}

	slog.Info("created build", slog.Int64("build_id", buildId))

	time.Sleep(5 * time.Second)
	err = w.buildRepo.UpdateStatus(w.ctx, buildId, "in_progress")
	if err != nil {
		slog.Error("failed to update build status", slog.Any("error", err))
		return
	}

	slog.Debug("build in progress", slog.Int64("build_id", buildId))

	time.Sleep(5 * time.Second)

	// random failure or success, 50% chance of failure
	conclusion := "success"
	if rand.Intn(2) == 0 {
		conclusion = "failure"
	}

	err = w.buildRepo.SetConclusion(w.ctx, buildId, conclusion)
	if err != nil {
		slog.Error("failed to set build conclusion", slog.Any("error", err))
		return
	}

	slog.Debug("build finished", slog.Int64("build_id", buildId), slog.String("conclusion", conclusion))
}
