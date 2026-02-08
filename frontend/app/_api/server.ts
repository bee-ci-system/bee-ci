import { serverFetch } from '../_utils/server-fetch';

interface GetUserDto {
  id: number;
  name: string;
}

export async function getUserServer(): Promise<GetUserDto> {
  const res = await serverFetch('/user');

  if (!res) {
    throw new Error('fetching failed on /user');
  }

  if (!res.ok) {
    throw new Error(await res.text());
  }

  return (await res.json()) as GetUserDto;
}
