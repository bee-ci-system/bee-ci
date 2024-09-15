import { getPipelineInfoServer, getPipelineLogsServer } from './_api/server';
import { PipelineInfoCard } from './_components/pipeline-info-card';
import { TerminalWindow } from './_components/terminal';

const PipelinePage = async ({ params }: { params: { id: string } }) => {
  const [pipelineInfo, pipelineLogs] = await Promise.all([
    getPipelineInfoServer(params.id),
    getPipelineLogsServer(params.id),
  ]);

  return (
    <div className='mb-4 mt-4 flex h-[90%] max-w-[1800px] flex-col gap-4 md:mt-0 md:flex-row'>
      <div className='mx-4 flex h-fit flex-grow md:h-[90vh] md:w-1/3'>
        <PipelineInfoCard pipeline={pipelineInfo} />
      </div>
      <div className='mx-4 flex h-[500px] flex-grow md:h-[90vh] md:w-2/3'>
        <TerminalWindow
          pipelineId={pipelineInfo.id}
          initialLogs={pipelineLogs}
        />
      </div>
    </div>
  );
};

export default PipelinePage;
