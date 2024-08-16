import { ChevronLeftIcon, ChevronRightIcon } from 'lucide-react';
import Link from 'next/link';
import { getPreviousNext } from '../_utils/markdown';

export default function Pagination({ pathname }: { pathname: string }) {
  const paginationData = getPreviousNext(pathname);

  return (
    <div className='flex items-center justify-between py-5 sm:py-7'>
      <div>
        {paginationData.prev && (
          <Link
            className='flex items-center gap-2 px-1 text-sm no-underline'
            href={`${paginationData.prev.href}`}
          >
            <ChevronLeftIcon className='h-[1.1rem] w-[1.1rem]' />
            <p>{paginationData.prev.title}</p>
          </Link>
        )}
      </div>
      <div>
        {paginationData.next && (
          <Link
            className='flex items-center gap-2 px-1 text-sm no-underline'
            href={`${paginationData.next.href}`}
          >
            <p>{paginationData.next.title}</p>
            <ChevronRightIcon className='h-[1.1rem] w-[1.1rem]' />
          </Link>
        )}
      </div>
    </div>
  );
}
