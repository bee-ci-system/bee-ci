import {
  GetMyRepositoriesDataDto,
  GetMyRepositoriesParams,
  Repository,
} from '@/app/_types/my-repositories';

const repositoriesData: Repository[] = [
  { id: '1', name: 'kacaleksandra/project1', dateOfLastUpdate: '2023-01-01' },
  { id: '2', name: 'kacaleksandra/project2', dateOfLastUpdate: '2023-01-02' },
  {
    id: '3',
    name: 'kacaleksandra/my-awesome-project',
    dateOfLastUpdate: '2023-01-03',
  },
  { id: '4', name: 'kacaleksandra/cool-app', dateOfLastUpdate: '2023-01-04' },
  { id: '5', name: 'kacaleksandra/test-repo', dateOfLastUpdate: '2023-01-05' },
  {
    id: '6',
    name: 'kacaleksandra/sample-project',
    dateOfLastUpdate: '2023-01-06',
  },
  { id: '7', name: 'kacaleksandra/demo-repo', dateOfLastUpdate: '2023-01-07' },
  { id: '8', name: 'kacaleksandra/repo-eight', dateOfLastUpdate: '2023-01-08' },
  { id: '9', name: 'kacaleksandra/repo-nine', dateOfLastUpdate: '2023-01-09' },
  { id: '10', name: 'kacaleksandra/repo-ten', dateOfLastUpdate: '2023-01-10' },
  {
    id: '11',
    name: 'kacaleksandra/another-repo',
    dateOfLastUpdate: '2023-01-11',
  },
  { id: '12', name: 'kacaleksandra/some-repo', dateOfLastUpdate: '2023-01-12' },
  { id: '13', name: 'kacaleksandra/cool-repo', dateOfLastUpdate: '2023-01-13' },
  {
    id: '14',
    name: 'kacaleksandra/fun-project',
    dateOfLastUpdate: '2023-01-14',
  },
  {
    id: '15',
    name: 'kacaleksandra/unique-repo',
    dateOfLastUpdate: '2023-01-15',
  },
  {
    id: '16',
    name: 'kacaleksandra/repo-sixteen',
    dateOfLastUpdate: '2023-01-16',
  },
  {
    id: '17',
    name: 'kacaleksandra/repo-seventeen',
    dateOfLastUpdate: '2023-01-17',
  },
  {
    id: '18',
    name: 'kacaleksandra/repo-eighteen',
    dateOfLastUpdate: '2023-01-18',
  },
  {
    id: '19',
    name: 'kacaleksandra/repo-nineteen',
    dateOfLastUpdate: '2023-01-19',
  },
  {
    id: '20',
    name: 'kacaleksandra/repo-twenty',
    dateOfLastUpdate: '2023-01-20',
  },
  {
    id: '95a457e8-e9ff-4834-b698-2b6ab3c06715',
    name: 'kacaleksandra/bee-ci',
    dateOfLastUpdate: '2021-10-01',
  },
  {
    id: '95a457e8-e9ff-4834-b698-2b6ab3c06715',
    name: 'kacaleksandra/flashwise',
    dateOfLastUpdate: '2021-10-01',
  },
  {
    id: '95a457e8-e9ff-4834-b698-2b6ab3c06715',
    name: 'kacaleksandra/taskshare',
    dateOfLastUpdate: '2021-10-01',
  },
];

const filterAndPaginateRepositories = (
  search: string,
  page: number,
  pageSize: number,
): GetMyRepositoriesDataDto => {
  const filteredRepositories = repositoriesData.filter((repo) =>
    repo.name.toLowerCase().includes(search.toLowerCase()),
  );

  const totalRepositories = filteredRepositories.length;
  const totalPages = Math.ceil(totalRepositories / pageSize);
  const startIndex = (page - 1) * pageSize;
  const paginatedRepositories = filteredRepositories.slice(
    startIndex,
    startIndex + pageSize,
  );

  return {
    repositories: paginatedRepositories,
    totalRepositories,
    totalPages,
    currentPage: page,
  };
};

export const getMyRepositoriesData = async (
  params: GetMyRepositoriesParams,
): Promise<GetMyRepositoriesDataDto> => {
  const { currentPage, search } = params;
  const pageSize = 5;

  return filterAndPaginateRepositories(search, currentPage, pageSize);
};
