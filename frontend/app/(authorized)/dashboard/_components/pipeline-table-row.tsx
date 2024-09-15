import { buttonVariants } from '@/app/_components/button';
import { TableCell, TableRow } from '@/app/_components/table';
import { PipelineDashboardData } from '@/app/_types/dashboard';
import { cn } from '@/app/_utils/cn';
import { PipelineStatusToIcon } from '@/app/_utils/pipeline-status-to-icon';
import { routeGenerators } from '@/app/_utils/routes';
import { ArrowRight } from 'lucide-react';
import Link from 'next/link';

const PipelineTableRow = ({
  pipeline,
}: {
  pipeline: PipelineDashboardData;
}) => (
  <TableRow className='w-full'>
    <TableCell className='w-[75%]'>
      <div className='font-medium'>{pipeline.repositoryName}</div>
      <div className='text-sm text-muted-foreground md:inline'>
        {pipeline.commitName}
      </div>
    </TableCell>
    <TableCell className='flex items-center justify-between gap-4 px-0'>
      <p>{PipelineStatusToIcon(pipeline.status)}</p>
      <Link
        href={routeGenerators.pipeline(pipeline.id)}
        aria-label='open info about pipeline'
        className={cn(buttonVariants({ size: 'icon', variant: 'ghost' }))}
      >
        <ArrowRight
          width={18}
          height={18}
          strokeWidth={2}
          className='dark:text-white'
        />
      </Link>
    </TableCell>
  </TableRow>
);

export { PipelineTableRow };
