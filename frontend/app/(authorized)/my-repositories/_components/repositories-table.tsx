'use client';

import { Button } from '@/app/_components/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/app/_components/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/app/_components/table';
import { GetMyRepositoriesDataDto } from '@/app/_types/my-repositories';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeftIcon, ArrowRightIcon } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { getMyRepositoriesData } from '../_api/server';
import { Search } from './search';

const RepositoriesTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [search, setSearch] = useState('');

  const router = useRouter();

  const { data } = useQuery<GetMyRepositoriesDataDto>({
    queryKey: ['repositories', currentPage, search],
    queryFn: () =>
      getMyRepositoriesData({
        currentPage,
        search,
      }),
  });

  return (
    <Card className='flex flex-grow flex-col'>
      <CardHeader>
        <CardTitle>My repositories</CardTitle>
        <CardDescription>
          Here you can find all of your repositories and search for specific
          ones.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Search
          searchValue={search}
          handleSearchChange={(value) => {
            setSearch(value);
            setCurrentPage(1);
          }}
        />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className='w-2/3'>Name</TableHead>
              <TableHead>Last updated at</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data?.repositories.map((repository) => (
              <TableRow
                key={repository.id}
                onClick={() => router.push(`/repository/${repository.id}`)}
              >
                <TableCell className='font-medium'>{repository.name}</TableCell>
                <TableCell>{repository.dateOfLastUpdate}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
      {data && data.totalPages > 1 && (
        <CardFooter className='flex flex-grow flex-col gap-2'>
          <div className='flex gap-8'>
            <Button
              variant='outline'
              disabled={currentPage === 1}
              onClick={() => setCurrentPage((prev) => prev - 1)}
            >
              <ArrowLeftIcon />
            </Button>
            <Button
              variant='outline'
              disabled={currentPage + 1 > data?.totalPages}
              onClick={() => setCurrentPage((prev) => prev + 1)}
            >
              <ArrowRightIcon />
            </Button>
          </div>
          <div className='justify-right w-full text-right text-xs text-muted-foreground'>
            <strong>{data?.totalRepositories}</strong> repositories found
          </div>
        </CardFooter>
      )}
    </Card>
  );
};

export { RepositoriesTable };
