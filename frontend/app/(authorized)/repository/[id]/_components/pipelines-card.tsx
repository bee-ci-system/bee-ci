import { Card, CardContent, CardHeader } from '@/app/_components/card';
import { ScrollArea } from '@/app/_components/scroll-area';
import { Pipeline } from '@/app/_types/pipeline';
import { PipelineStatusToIcon } from '@/app/_utils/pipeline-status-to-icon';
import { routeGenerators } from '@/app/_utils/routes';
import { format } from 'date-fns';
import Link from 'next/link';

const PipelinesCard = ({ pipelines }: { pipelines: Pipeline[] }) => (
  <Card className='flex w-full flex-col rounded-none rounded-l-lg'>
    <CardHeader>
      <h2 className='text-beeci-yellow-500 dark:text-beeci-yellow-400'>
        Builds
      </h2>
    </CardHeader>
    <CardContent>
      {pipelines.length === 0 ? (
        <div className='flex h-[420px] items-center justify-center md:h-[80vh]'>
          <p className='text-muted-foreground'>You don't have builds yet.</p>
        </div>
      ) : (
        <ScrollArea className='h-[420px] md:h-[80vh]'>
          {pipelines.map((pipeline) => (
            <Link
              href={routeGenerators.pipeline(pipeline.id)}
              key={pipeline.id}
              className='mb-4 mr-4 flex cursor-pointer items-center gap-2 rounded-md border p-4 shadow-md hover:bg-primary-foreground'
            >
              <div className='mr-2 w-3/5'>
                <h3 className='text-base font-medium'>{pipeline.commitName}</h3>
                <p className='text-sm text-muted-foreground'>
                  {format(pipeline.startDate, 'HH:mm - dd MMM yyyy')}
                </p>
              </div>
              <div className='flex w-2/5 justify-end'>
                <p className='flex items-center gap-2'>
                  <span className='text-sm'>{pipeline.status}</span>
                  {PipelineStatusToIcon(pipeline.status)}
                </p>
              </div>
            </Link>
          ))}
        </ScrollArea>
      )}
    </CardContent>
  </Card>
);

export { PipelinesCard };
