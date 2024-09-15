// Package worker implements a thing that accepts jobs and schedules them for
// execution.
//
// After the job is added, it is picked up by one of the active executors.
//
// It spawns a single goroutine per new job.
package worker

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/bee-ci/bee-ci-system/data"
)

type Worker struct {
	ctx       context.Context
	buildRepo data.BuildRepo
}

// New creates a new [Worker].
//
// The worker can be scheduled with [Add] method. All jobs will
func New(ctx context.Context, buildRepo data.BuildRepo) *Worker {
	return &Worker{
		ctx:       ctx,
		buildRepo: buildRepo,
	}
}

// Add schedules a new job for execution.
//
// The job will be canceled when the context that was passed to [New] is
// canceled.
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

func SleepContext(ctx context.Context, d time.Duration) {
	select {
	case <-time.After(d):
	case <-ctx.Done():
	}
}
