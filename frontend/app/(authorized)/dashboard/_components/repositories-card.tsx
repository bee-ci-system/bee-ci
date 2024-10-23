import { Button, buttonVariants } from '@/app/_components/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/app/_components/card';
import { Repository } from '@/app/_types/my-repositories';
import { cn } from '@/app/_utils/cn';
import { routeGenerators, routes } from '@/app/_utils/routes';
import { ArrowRight, ArrowUpRight } from 'lucide-react';
import Link from 'next/link';

const MyRepositoriesCard = ({
  repositories,
}: {
  repositories: Repository[];
}) => (
  <Card className='border-beeci-yellow-600'>
    <CardHeader className='flex flex-row items-center'>
      <div className='grid gap-2'>
        <CardTitle>Repositories</CardTitle>
        <CardDescription>A list of all your repositories</CardDescription>
      </div>
      <Button
        asChild
        size='sm'
        className='ml-auto gap-1 bg-beeci-yellow-500 dark:bg-beeci-yellow-600'
      >
        <Link
          href={
            repositories.length === 0
              ? 'https://github.com/apps/bee-ci-system'
              : routes.MY_REPOSITORIES
          }
        >
          {repositories.length === 0 ? 'Configure' : 'View All'}
          <ArrowUpRight className='h-4 w-4' />
        </Link>
      </Button>
    </CardHeader>
    <CardContent className='grid gap-8'>
      {repositories.length === 0 && (
        <p className='mb-4 h-full text-center text-sm text-muted-foreground'>
          No repositories found
        </p>
      )}
      {repositories.map((repository) => (
        <div className='flex items-center gap-4' key={repository.id}>
          <div className='flex w-full justify-between'>
            <div className='flex flex-col gap-1'>
              <p className='text-sm font-medium leading-none'>
                {repository.name}
              </p>
              <p className='text-sm text-muted-foreground'>
                {repository.dateOfLastUpdate || 'no updates'}
              </p>
            </div>
            <Link
              href={routeGenerators.repository(repository.id)}
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
