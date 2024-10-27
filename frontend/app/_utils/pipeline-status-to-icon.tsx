import { CircleCheck, CircleX, LoaderCircle } from 'lucide-react';
import { PipelineDashboardData } from '../_types/dashboard';
import {
  Pipeline,
  PipelineConclusion,
  PipelineStatus,
} from '../_types/pipeline';

const PipelineStatusToIcon = (pipeline: Pipeline | PipelineDashboardData) => {
  const state =
    pipeline.status === PipelineStatus.COMPLETED && pipeline.conclusion
      ? pipeline.conclusion
      : pipeline.status;
  switch (state) {
    case PipelineConclusion.SUCCESS:
      return (
        <CircleCheck width={24} height={24} className='text-emerald-500' />
      );
    case PipelineConclusion.FAILURE:
      return <CircleX width={24} height={24} className='text-red-500' />;
    case PipelineStatus.IN_PROGRESS:
      return (
        <svg
          className='h-6 w-6 animate-spin text-yellow-500'
          viewBox='0 0 24 24'
        >
          <LoaderCircle />
        </svg>
      );
    case PipelineStatus.QUEUED:
      return (
        <svg
          className='h-6 w-6 animate-pulse text-blue-500'
          viewBox='0 0 24 24'
        >
          <LoaderCircle />
        </svg>
      );
  }
};

export { PipelineStatusToIcon };
