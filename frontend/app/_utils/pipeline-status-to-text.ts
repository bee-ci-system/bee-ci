import { PipelineDashboardData } from '../_types/dashboard';
import { Pipeline, PipelineStatus } from '../_types/pipeline';

const PipelineStatusToText = (pipeline: Pipeline | PipelineDashboardData) => {
  const state =
    pipeline.status === PipelineStatus.COMPLETED && pipeline.conclusion
      ? pipeline.conclusion
      : pipeline.status;

  switch (state) {
    case PipelineStatus.IN_PROGRESS:
      return 'in progress';
    default:
      return state;
  }
};

export { PipelineStatusToText };
