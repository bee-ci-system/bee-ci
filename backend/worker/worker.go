// Package worker implements a thing that accepts jobs and schedules them for
// execution.
//
// After the job is added, it is picked up by one of the active executors.
//
// It spawns a single goroutine per new job.
package worker

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/bee-ci/bee-ci-system/internal/data"
)

type Worker struct {
	ctx       context.Context
	buildRepo data.BuildRepo
	logger    *slog.Logger
}

// New creates a new [Worker].
//
// The worker can be scheduled with [Worker.Add] method.
func New(ctx context.Context, buildRepo data.BuildRepo) *Worker {
	return &Worker{
		logger:    slog.Default(), // TODO: add some "subsystem name" to this logger
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
		w.logger.Error("failed to create build", slog.Any("error", err))
		// TODO: handle error in a better way – update status on GitHub
		return
	}

	w.logger.Debug("created build", slog.Int64("build_id", buildId))

	err = SleepContext(w.ctx, 5*time.Second)
	if err != nil {
		w.logger.Error("job processing aborted", slog.Any("error", err))
		return
	}

	err = w.buildRepo.UpdateStatus(w.ctx, buildId, "in_progress")
	if err != nil {
		w.logger.Error("failed to update build status", slog.Any("error", err))
		return
	}

	w.logger.Debug("build in progress", slog.Int64("build_id", buildId))

	err = SleepContext(w.ctx, 5*time.Second)
	if err != nil {
		w.logger.Error("job processing aborted", slog.Any("error", err))
	}

	// random failure or success, 50% chance of failure
	conclusion := "success"
	if rand.Intn(2) == 0 {
		conclusion = "failure"
	}

	err = w.buildRepo.SetConclusion(w.ctx, buildId, conclusion)
	if err != nil {
		w.logger.Error("failed to set build conclusion", slog.Any("error", err))
		return
	}

	w.logger.Debug("build finished", slog.Int64("build_id", buildId), slog.String("conclusion", conclusion))
}

// TODO: Create an issue for this func to be added to Go stdlib.
func SleepContext(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return fmt.Errorf("sleep was aborted because context is done before duration passed: %w", ctx.Err())
	}
}
