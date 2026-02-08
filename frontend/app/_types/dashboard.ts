import { Repository } from './my-repositories';
import { PipelineConclusion, PipelineStatus } from './pipeline';

export interface PipelineDashboardData {
  id: string;
  repositoryName: string;
  commitName: string;
  status: PipelineStatus;
  conclusion: PipelineConclusion | null;
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
