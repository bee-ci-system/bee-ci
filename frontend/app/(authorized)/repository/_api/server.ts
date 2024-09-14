import { PipelineStatus } from '@/app/_types/pipeline';
import { GetRepositoryDto } from '@/app/_types/repository';

export const getRepositoryData = async ({
  id,
}: {
  id: string;
}): Promise<GetRepositoryDto> => {
  return {
    id,
    name: 'web-app-frontend',
    description: 'Frontend repository for the web application.',
    url: 'https://github.com/company/web-app-frontend',
    dateOfLastUpdate: '2024-09-13',
    pipelines: [
      {
        id: 'pipeline-001',
        repositoryName: 'web-app-frontend',
        commitName: 'feature/add-user-authentication',
        status: PipelineStatus.IN_PROGRESS,
        startDate: '2024-09-13T10:15:00',
      },
      {
        id: 'pipeline-002',
        repositoryName: 'web-app-frontend',
        commitName: 'bugfix/fix-login-redirect',
        status: PipelineStatus.SUCCESS,
        startDate: '2024-09-12T09:45:00',
      },
      {
        id: 'pipeline-003',
        repositoryName: 'web-app-frontend',
        commitName: 'chore/update-dependencies',
        status: PipelineStatus.SUCCESS,
        startDate: '2024-09-13T11:00:00',
      },
      {
        id: 'pipeline-004',
        repositoryName: 'web-app-frontend',
        commitName: 'feature/add-dark-mode',
        status: PipelineStatus.FAILURE,
        startDate: '2024-09-13T11:30:00',
      },
      {
        id: 'pipeline-005',
        repositoryName: 'web-app-frontend',
        commitName: 'feature/add-user-profile',
        status: PipelineStatus.SUCCESS,
        startDate: '2024-09-13T12:00:00',
      },
      {
        id: 'pipeline-006',
        repositoryName: 'web-app-frontend',
        commitName: 'feature/add-user-settings',
        status: PipelineStatus.FAILURE,
        startDate: '2024-09-13T12:30:00',
      },
      {
        id: 'pipeline-007',
        repositoryName: 'web-app-frontend',
        commitName: 'feature/add-user-activity',
        status: PipelineStatus.SUCCESS,
        startDate: '2024-09-13T13:00:00',
      },
      {
        id: 'pipeline-008',
        repositoryName: 'web-app-frontend',
        commitName: 'feature/add-user-projects',
        status: PipelineStatus.SUCCESS,
        startDate: '2024-09-13T13:30:00',
      },
      {
        id: 'pipeline-009',
        repositoryName: 'web-app-frontend',
        commitName: 'feature/add-user-teams',
        status: PipelineStatus.FAILURE,
        startDate: '2024-09-13T14:00:00',
      },
      {
        id: 'pipeline-010',
        repositoryName: 'web-app-frontend',
        commitName: 'feature/add-user-organizations',
        status: PipelineStatus.FAILURE,
        startDate: '2024-09-13T14:30:00',
      },
      {
        id: 'pipeline-011',
        repositoryName: 'web-app-frontend',
        commitName: 'feature/add-user-roles',
        status: PipelineStatus.FAILURE,
        startDate: '2024-09-13T15:00:00',
      },
    ],
  };
};
