import { GetDashboardDataDto } from '@/app/_types/dashboard';
import { PipelineStatus } from '@/app/_types/pipeline';

export const getDashboardData = async (): Promise<GetDashboardDataDto> => {
  return {
    stats: {
      totalPipelines: 13,
      successfulPipelines: 3,
      unsuccessfulPipelines: 10,
    },
    repositories: [
      {
        id: '95a457e8-e9ff-4834-b698-2b6ab3c06715',
        name: 'kacaleksandra/bee-ci',
        dateOfLastUpdate: '2021-10-01',
      },
    ],
    pipelines: [
      {
        id: 'd5bab913-49ef-4058-85ff-3f0dfd1e7023',
        repositoryName: 'kacaleksandra/bee-ci',
        commitName: 'feature/BEECI-528/added-dashboard',
        status: PipelineStatus.FAILURE,
      },
      {
        id: '2018c1fc-8ae8-4bd3-8322-50f8d6bdcfe6',
        repositoryName: 'kacaleksandra/flashwise',
        commitName: 'feature/flashwise/added-new-types',
        status: PipelineStatus.SUCCESS,
      },
      {
        id: '2fbdf5d9-1495-4c61-979b-7a9f4efa6d41',
        repositoryName: 'kacaleksandra/taskshare',
        commitName: 'fix/task-share/fixed-api-communication',
        status: PipelineStatus.IN_PROGRESS,
      },
    ],
  };
};
