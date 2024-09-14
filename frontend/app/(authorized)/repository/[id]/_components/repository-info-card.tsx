import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
} from '@/app/_components/card';
import { routes } from '@/app/_utils/routes';
import { format } from 'date-fns';
import { CircleChevronLeft, Clock3, SquareArrowOutUpRight } from 'lucide-react';
import Link from 'next/link';

const RepositoryInfoCard = ({
  name,
  description,
  dateOfLastUpdate,
  url,
}: {
  name: string;
  description: string;
  dateOfLastUpdate: string;
  url: string;
}) => (
  <Card className='flex w-full flex-col pr-6'>
    <CardHeader>
      <h2 className='text-beeci-yellow-500 dark:text-beeci-yellow-400'>
        {name}
      </h2>
      <CardDescription>{description}</CardDescription>
    </CardHeader>
    <CardContent className='flex flex-grow flex-col gap-8 text-sm text-foreground'>
      <div>
        <p className='w-full leading-loose'>
          <span className='flex items-center gap-2'>
            <Clock3 className='size-4' /> Date of last update:
          </span>
          <span className='block w-full text-right text-base'>
            {format(dateOfLastUpdate, 'HH:mm - dd MMM yyyy')}
          </span>
        </p>
      </div>
      <div>
        <p className='leading-loose'>
          <span className='flex w-full items-center gap-2'>
            <SquareArrowOutUpRight className='size-4' /> Url:
          </span>
          <Link
            className='ml-4 block w-full break-words text-right text-base underline'
            href={url}
          >
            {url}
          </Link>
        </p>
      </div>
    </CardContent>
    <CardFooter className='my-4'>
      <Link href={routes.MY_REPOSITORIES}>
        <CircleChevronLeft className='text-beeci-yellow-500 dark:text-beeci-yellow-400' />
      </Link>
    </CardFooter>
  </Card>
);

export { RepositoryInfoCard };
