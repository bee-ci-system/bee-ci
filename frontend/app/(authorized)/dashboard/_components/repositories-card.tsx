import { Button, buttonVariants } from '@/app/_components/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/app/_components/card';
import { RepositoriesDashboardData } from '@/app/_types/dashboard';
import { cn } from '@/app/_utils/cn';
import { ArrowRight, ArrowUpRight } from 'lucide-react';
import Link from 'next/link';

const MyRepositoriesCard = ({
  repositories,
}: {
  repositories: RepositoriesDashboardData[];
}) => (
  <Card className='border-beeci-yellow-600' x-chunk='dashboard-01-chunk-3'>
    <CardHeader className='flex flex-row items-center'>
      <div className='grid gap-2'>
        <CardTitle>Repositories</CardTitle>
        <CardDescription>A list of all your repositories</CardDescription>
      </div>
      <Button asChild size='sm' className='ml-auto gap-1 bg-beeci-yellow-600'>
        <Link href='#'>
          View All
          <ArrowUpRight className='h-4 w-4' />
        </Link>
      </Button>
    </CardHeader>
    <CardContent className='grid gap-8'>
      {repositories.map((repository) => (
        <div className='flex items-center gap-4' key={repository.id}>
          <div className='flex w-full justify-between'>
            <div className='grid gap-1'>
              <p className='text-sm font-medium leading-none'>
                {repository.name}
              </p>
              <p className='text-sm text-muted-foreground'>
                {repository.dateOfLastUpdate}
              </p>
            </div>
            <Link
              href='#'
              aria-label='open info about pipeline'
              className={cn(buttonVariants({ size: 'icon', variant: 'ghost' }))}
            >
              <ArrowRight
                width={24}
                height={24}
                strokeWidth={2}
                className='dark:text-white'
              />
            </Link>
          </div>
        </div>
      ))}
    </CardContent>
  </Card>
);

export { MyRepositoriesCard };
