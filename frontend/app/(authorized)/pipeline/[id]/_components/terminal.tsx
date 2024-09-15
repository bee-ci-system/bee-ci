'use client';

import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/app/_components/tooltip';
import { TooltipProvider } from '@radix-ui/react-tooltip';
import { keepPreviousData, useQuery } from '@tanstack/react-query';
import { useEffect, useRef, useState } from 'react';
import {
  getAllPipelineLogsClient,
  getPipelineLogsClient,
} from '../_api/client';

const TerminalWindow = ({
  pipelineId,
  initialLogs,
}: {
  pipelineId: string;
  initialLogs: string;
}) => {
  const [fetchAllLogs, setFetchAllLogs] = useState(false);

  const { data: logsData } = useQuery<string>({
    queryKey: ['pipelineLogs', pipelineId, fetchAllLogs],
    queryFn: () =>
      fetchAllLogs
        ? getAllPipelineLogsClient(pipelineId)
        : getPipelineLogsClient(pipelineId),
    placeholderData: keepPreviousData,
    initialData: initialLogs,
    refetchInterval: fetchAllLogs ? undefined : 2000,
    enabled: true,
  });

  const lines = logsData.split('\n');

  const logEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (logEndRef.current) {
      logEndRef.current.scrollIntoView({ behavior: 'instant' });
    }
  }, [logsData]);

  return (
    <div className='h-full w-full'>
      <div className='coding inverse-toggle h-full rounded-lg bg-slate-100 px-5 pb-6 pt-4 font-mono text-sm leading-normal text-foreground subpixel-antialiased shadow-lg dark:bg-gray-800'>
        <div className='top mb-2 flex'>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <div
                  className='h-3 w-3 cursor-pointer rounded-full bg-orange-300'
                  role='button'
                  onClick={() => setFetchAllLogs(true)}
                />
              </TooltipTrigger>
              <TooltipContent side='bottom'>
                <p>show all logs</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
        {logsData === '' ? (
          <div className='flex h-full items-center justify-center'>
            <p className='text-muted-foreground'>No logs</p>
          </div>
        ) : (
          <div className='h-full overflow-y-scroll'>
            {lines.map((line, index) => (
              <div key={index} className='mb-1 flex'>
                <span className='text-beeci-yellow-500 dark:text-beeci-yellow-400'>
                  logs:~$
                </span>
                <p className='typing flex-1 items-center pl-2'>{line}</p>
              </div>
            ))}
            <div ref={logEndRef} />
          </div>
        )}
      </div>
    </div>
  );
};

export { TerminalWindow };
