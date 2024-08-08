// Package worker implements a worker that executes jobs.
//
// It spawns a single goroutine per new job.
package worker

import (
	"context"
	"time"

	"github.com/bartekpacia/ghapp/data"
)

type Worker interface {
	Add(build data.NewBuildRequest)
}

type worker struct {
	ctx       context.Context
	buildRepo data.BuildRepo
}

func (w worker) Add(build data.NewBuildRequest) {
	go w.job(build)
}

func (w worker) job(build data.NewBuildRequest) {
	// Do some work
	buildId, err := w.buildRepo.Create(w.ctx, build)
	if err != nil {
		// TODO: handle error
	}

	// job is queued

	time.Sleep(5 * time.Second)
	w.buildRepo.Update(w.ctx, buildId, data.StatusRunning)

	time.Sleep(5 * time.Second)

	// random failure or success, 50% chance of failure
	if time.Now().UnixNano()%2 == 0 {
		w.buildRepo.Update(w.ctx, buildId, data.StatusFailed)
	} else {
		w.buildRepo.Update(w.ctx, buildId, data.StatusSuccess)
	}
}

func New(ctx context.Context, buildRepo data.BuildRepo) Worker {
	return &worker{
		buildRepo: buildRepo,
	}
}
