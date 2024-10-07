import { GetDashboardDataDto } from '@/app/_types/dashboard';
import { serverFetch } from '@/app/_utils/server-fetch';

export const getDashboardDataServer =
  async (): Promise<GetDashboardDataDto> => {
    const res = await serverFetch('/dashboard');

    if (!res) {
      throw new Error('fetching failed on /dashboard');
    }

    if (!res.ok) {
      throw new Error(await res.text());
    }

    return (await res.json()) as GetDashboardDataDto;
  };
