import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
} from '@/app/_components/card';
import { Pipeline } from '@/app/_types/pipeline';
import { PipelineStatusToIcon } from '@/app/_utils/pipeline-status-to-icon';
import { PipelineStatusToText } from '@/app/_utils/pipeline-status-to-text';
import { routeGenerators } from '@/app/_utils/routes';
import { format, formatDistance } from 'date-fns';
import {
  BookCheck,
  Check,
  CircleChevronLeft,
  Clock2,
  Play,
} from 'lucide-react';
import Link from 'next/link';

const PipelineInfoCard = ({ pipeline }: { pipeline: Pipeline }) => (
  <Card className='flex w-full flex-col'>
    <CardHeader className='mb-8 border-b'>
      <h2 className='text-beeci-yellow-500 dark:text-beeci-yellow-400'>
        {pipeline.commitName}
      </h2>
      <CardDescription>from repo: {pipeline.repositoryName}</CardDescription>
    </CardHeader>
    <CardContent className='mr-6 flex flex-grow flex-col gap-8 text-sm text-foreground'>
      <div>
        <p className='w-full leading-loose'>
          <span className='flex items-center gap-2'>
            <BookCheck className='size-4' /> Status
          </span>
          <span className='flex w-full items-center justify-end gap-2'>
            {PipelineStatusToIcon(pipeline)}
            <span className='text-sm font-medium'>
              {PipelineStatusToText(pipeline)}
            </span>
          </span>
        </p>
      </div>
      <div>
        <p className='leading-loose'>
          <span className='flex w-full items-center gap-2'>
            <Play className='size-4' /> Start date:
          </span>
          <span className='block w-full text-right text-base'>
            {format(new Date(pipeline.startDate), 'HH:mm - dd MMM yyyy')}
          </span>
        </p>
      </div>
      <div>
        <p className='leading-loose'>
          <span className='flex w-full items-center gap-2'>
            <Check className='size-4' /> End date:
          </span>
          <span className='block w-full text-right text-base'>
            {pipeline.endDate
              ? format(pipeline.endDate, 'HH:mm - dd MMM yyyy')
              : '-'}
          </span>
        </p>
      </div>
      <div>
        <p className='leading-loose'>
          <span className='flex w-full items-center gap-2'>
            <Clock2 className='size-4' /> Time:
          </span>
          <span className='block w-full text-right text-base'>
            {pipeline.endDate
              ? formatDistance(pipeline.endDate, pipeline.startDate)
              : '-'}
          </span>
        </p>
      </div>
    </CardContent>
    <CardFooter className='my-4'>
      <Link href={routeGenerators.repository(pipeline.repositoryId)}>
        <CircleChevronLeft className='text-beeci-yellow-500 dark:text-beeci-yellow-400' />
      </Link>
    </CardFooter>
  </Card>
);

export { PipelineInfoCard };
