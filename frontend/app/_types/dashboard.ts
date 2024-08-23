import { PipelineStatus } from './pipeline';

export interface PipelineDashboardData {
  id: string;
  repositoryName: string;
  commitName: string;
  status: PipelineStatus;
}

export interface RepositoriesDashboardData {
  id: string;
  name: string;
  dateOfLastUpdate: string;
}

export interface GetDashboardDataDto {
  stats: {
    totalPipelines: number;
    successfulPipelines: number;
    unsuccessfulPipelines: number;
  };
  repositories: RepositoriesDashboardData[];
  pipelines: PipelineDashboardData[];
}
