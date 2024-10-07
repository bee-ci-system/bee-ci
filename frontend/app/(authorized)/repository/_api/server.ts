import { GetRepositoryDto } from '@/app/_types/repository';
import { serverFetch } from '@/app/_utils/server-fetch';

export const getRepositoryDataServer = async ({
  id,
}: {
  id: string;
}): Promise<GetRepositoryDto> => {
  const res = await serverFetch(`/repositories/${id}`);

  if (!res) {
    throw new Error(`fetching failed on /repository/${id}`);
  }

  if (!res.ok) {
    throw new Error(await res.text());
  }

  return (await res.json()) as GetRepositoryDto;
};
