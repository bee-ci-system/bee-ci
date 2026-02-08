import { Pipeline } from './pipeline';

export interface GetRepositoryDto {
  id: string;
  name: string;
  description: string;
  url: string;
  dateOfLastUpdate: string;
  pipelines: Pipeline[];
}
