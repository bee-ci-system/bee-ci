export enum PipelineStatus {
  SUCCESS = 'success',
  FAILURE = 'failure',
  IN_PROGRESS = 'in progress',
  QUEUED = 'queued',
}

export interface Pipeline {
  id: string;
  repositoryName: string;
  repositoryId: string;
  commitName: string;
  status: PipelineStatus;
  startDate: string;
  endDate?: string;
}
