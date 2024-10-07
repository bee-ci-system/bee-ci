import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/app/_components/card';
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from '@/app/_components/table';
import { PipelineDashboardData } from '@/app/_types/dashboard';
import { PipelineTableRow } from './pipeline-table-row';

const PipelinesCard = ({
  pipelines,
}: {
  pipelines: PipelineDashboardData[];
}) => (
  <Card className='border-beeci-yellow-600 xl:col-span-2'>
    <CardHeader>
      <div className='grid gap-2'>
        <CardTitle>Pipelines</CardTitle>
        <CardDescription>Your last pipelines and their status</CardDescription>
      </div>
    </CardHeader>
    {pipelines.length === 0 ? (
      <p className='mb-4 h-full text-center text-sm text-muted-foreground'>
        No pipelines found
      </p>
    ) : (
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className='w-[75%]'>Pipeline</TableHead>
              <TableHead className='pl-0'>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {pipelines.map((pipeline) => (
              <PipelineTableRow pipeline={pipeline} key={pipeline.id} />
            ))}
          </TableBody>
        </Table>
      </CardContent>
    )}
  </Card>
);

export { PipelinesCard };
