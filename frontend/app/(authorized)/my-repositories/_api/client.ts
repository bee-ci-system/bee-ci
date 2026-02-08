import {
  GetMyRepositoriesDataDto,
  GetMyRepositoriesParams,
} from '@/app/_types/my-repositories';
import { clientFetch } from '@/app/_utils/client-fetch';

export const getMyRepositoriesDataClient = async (
  params: GetMyRepositoriesParams,
): Promise<GetMyRepositoriesDataDto> => {
  const { currentPage, search } = params;
  const pageSize = 5;

  const urlWithParams = new URLSearchParams();
  urlWithParams.append('currentPage', currentPage.toString());
  urlWithParams.append('pageSize', pageSize.toString());
  if (search) urlWithParams.append('search', search);

  const res = await clientFetch('/my-repositories?' + urlWithParams.toString());

  if (!res || !res.ok)
    return {
      repositories: [],
      totalRepositories: 0,
      totalPages: 1,
      currentPage: 0,
    };

  return (await res.json()) as GetMyRepositoriesDataDto;
};
