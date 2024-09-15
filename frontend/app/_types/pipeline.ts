export enum PipelineStatus {
  SUCCESS = 'success',
  FAILURE = 'failure',
  IN_PROGRESS = 'in progress',
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

export interface PipelineLogs {
  logs: string;
}
