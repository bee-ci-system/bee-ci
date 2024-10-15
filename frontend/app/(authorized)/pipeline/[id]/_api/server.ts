import { Pipeline } from '@/app/_types/pipeline';
import { serverFetch } from '@/app/_utils/server-fetch';

export const getPipelineInfoServer = async (
  pipelineId: string,
): Promise<Pipeline> => {
  const res = await serverFetch('/pipeline/' + pipelineId);

  if (!res) {
    throw new Error('fetching failed on /pipeline/' + pipelineId);
  }

  if (!res.ok) {
    throw new Error(await res.text());
  }

  const data = await res.json();

  return data as Pipeline;
};

export const getPipelineLogsServer = async (
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  pipelineId: string,
): Promise<string> => {
  const res = await serverFetch(`/pipeline/${pipelineId}/logs`);

  if (!res) {
    throw new Error(`fetching failed on /pipeline/${pipelineId}/logs`);
  }

  if (!res.ok) {
    throw new Error(await res.text());
  }

  const data = await res.text();

  return data;
};
