import { getRepositoryData } from '../_api/server';
import { PipelinesCard } from './_components/pipelines-card';
import { RepositoryInfoCard } from './_components/repository-info-card';

const RepositoryPage = async ({ params }: { params: { id: string } }) => {
  const repositoryData = await getRepositoryData({ id: params.id });

  return (
    <div className='mb-4 mt-4 flex h-[90%] max-w-[1800px] flex-col gap-4 md:mt-0 md:flex-row'>
      <div className='mx-4 flex h-fit flex-grow md:h-[90vh] md:w-1/3'>
        <RepositoryInfoCard
          name={repositoryData.name}
          description={repositoryData.description}
          dateOfLastUpdate={repositoryData.dateOfLastUpdate}
          url={repositoryData.url}
        />
      </div>
      <div className='mx-4 flex h-[500px] flex-grow md:mx-0 md:h-[90vh] md:w-2/3'>
        <PipelinesCard pipelines={repositoryData.pipelines} />
      </div>
    </div>
  );
};

export default RepositoryPage;
