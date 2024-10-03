package api

import "time"

type getMyRepositoriesDTO struct {
	Repositories      []repository `json:"repositories"`
	TotalRepositories int          `json:"totalRepositories"`
	TotalPages        int          `json:"totalPages"`
	CurrentPage       int          `json:"currentPage"`
}

type repository struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	DateOfLastUpdate time.Time `json:"string"`
}
