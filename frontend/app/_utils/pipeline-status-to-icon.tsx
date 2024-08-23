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
        <LoaderCircle width={24} height={24} className='text-yellow-500' />
      );
  }
};

export { PipelineStatusToIcon };
