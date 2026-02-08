import { clientFetch } from '@/app/_utils/client-fetch';

export const getPipelineLogsClient = async (
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  pipelineId: string,
): Promise<string> => {
  const res = await clientFetch(`/pipeline/${pipelineId}/logs`);

  if (!res || !res.ok) {
    return '';
  }

  return await res.text();
};
