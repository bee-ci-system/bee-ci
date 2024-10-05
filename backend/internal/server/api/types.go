package api

import (
	"github.com/bee-ci/bee-ci-system/internal/data"
	"strconv"
	"time"
)

type getMyRepositoriesParams struct {
	CurrentPage int    `json:"currentPage"`
	Search      string `json:"search"`
}

type getUserDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type getMyRepositoriesDTO struct {
	Repositories      []repository `json:"repositories"`
	TotalRepositories int          `json:"totalRepositories"`
	TotalPages        int          `json:"totalPages"`
	CurrentPage       int          `json:"currentPage"`
}

type repository struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	DateOfLastUpdate time.Time `json:"dateOfLastUpdate"`
}

type getDashboardDataDto struct {
	Stats        statsDTO                `json:"stats"`
	Repositories []repository            `json:"repositories"`
	Pipelines    []pipelineDashboardData `json:"pipelines"`
}

type statsDTO struct {
	TotalPipelines        int `json:"totalPipelines"`
	SuccessfulPipelines   int `json:"successfulPipelines"`
	UnsuccessfulPipelines int `json:"unsuccessfulPipelines"`
}

type pipelineDashboardData struct {
	ID             string `json:"id"`
	RepositoryName string `json:"repositoryName"`
	CommitName     string `json:"commitName"`
	Status         string `json:"status"`
}

func toRepositories(dbRepos []data.Repo) []repository {
	var repos []repository
	for _, repo := range dbRepos {
		repos = append(repos, repository{
			ID:               strconv.FormatInt(repo.ID, 10),
			Name:             repo.Name,
			DateOfLastUpdate: time.Date(2005, 04, 02, 21, 37, 0, 0, time.Local),
		})
	}
	return repos
}
