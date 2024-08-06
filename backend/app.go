package main

import "context"

type App struct {
	BuildService BuildService
}

type Build struct {
	RepoID uint64
	Commit string
}

type BuildService interface {
	Create(ctx context.Context, build Build) (id uint64, err error)
	GetByRepoID(ctx context.Context, repoID uint64) (builds []Build, err error)
}
