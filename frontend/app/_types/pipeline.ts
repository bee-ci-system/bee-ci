export enum PipelineStatus {
  COMPLETED = 'completed',
  IN_PROGRESS = 'in progress',
  QUEUED = 'queued',
}

export enum PipelineConclusion {
  FAILURE = 'failure',
  SUCCESS = 'success',
}

export interface Pipeline {
  id: string;
  repositoryName: string;
  repositoryId: string;
  commitName: string;
  status: PipelineStatus;
  conclusion: PipelineConclusion | null;
  startDate: string;
  endDate?: string;
}
