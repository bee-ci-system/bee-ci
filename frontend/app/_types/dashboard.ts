import { Repository } from './my-repositories';
import { PipelineStatus } from './pipeline';

export interface PipelineDashboardData {
  id: string;
  repositoryName: string;
  commitName: string;
  status: PipelineStatus;
}

export interface GetDashboardDataDto {
  stats: {
    totalPipelines: number;
    successfulPipelines: number;
    unsuccessfulPipelines: number;
  };
  repositories: Repository[];
  pipelines: PipelineDashboardData[];
}
