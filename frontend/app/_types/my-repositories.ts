export interface Repository {
  id: string;
  name: string;
  dateOfLastUpdate: string;
}

export interface GetMyRepositoriesParams {
  currentPage: number;
  search: string;
}

export interface GetMyRepositoriesDataDto {
  repositories: Repository[];
  totalRepositories: number;
  totalPages: number;
  currentPage: number;
}
