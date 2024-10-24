package api

import (
	"time"
)

type getMyRepositoriesParams struct {
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
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
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	DateOfLastUpdate *time.Time `json:"dateOfLastUpdate"`
}

type getDashboardDataDTO struct {
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

type getRepositoryDTO struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	URL              string     `json:"url"`
	DateOfLastUpdate *time.Time `json:"dateOfLastUpdate"`
	Pipelines        []pipeline `json:"pipelines"`
}

type pipeline struct {
	ID             string     `json:"id"`
	RepositoryName string     `json:"repositoryName"`
	RepositoryID   string     `json:"repositoryId"`
	CommitName     string     `json:"commitName"`
	Status         string     `json:"status"`
	StartDate      time.Time  `json:"startDate"`
	EndDate        *time.Time `json:"endDate"`
}
