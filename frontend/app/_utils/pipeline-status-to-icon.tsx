import { CircleCheck, CircleX, LoaderCircle } from 'lucide-react';
import { PipelineStatus } from '../_types/pipeline';

const PipelineStatusToIcon = (status: PipelineStatus) => {
  switch (status) {
    case PipelineStatus.SUCCESS:
      return (
        <CircleCheck width={24} height={24} className='text-emerald-500' />
      );
    case PipelineStatus.FAILURE:
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
