'use client';

import { FlagIcon } from 'lucide-react';
import Link from 'next/link';
import { Button } from './_components/button';
import { routes } from './_utils/routes';

export function ErrorPage() {
  return (
    <div className='mx-auto grid h-screen place-items-center px-8 text-center'>
      <div>
        <FlagIcon className='mx-auto h-20 w-20' />
        <h1 className='mt-10 !text-3xl !leading-snug md:!text-4xl'>
          Error 404 <br /> It looks like something went wrong.
        </h1>
        <h3 className='mx-auto mb-14 mt-8 text-[18px] font-normal text-gray-500 md:max-w-sm'>
          Don&apos;t worry, our team is already on it. Please try refreshing the
          page or come back later.
        </h3>
        <Link href={routes.DASHBOARD}>
          <Button>back dashboard</Button>
        </Link>
      </div>
    </div>
  );
}

export default ErrorPage;
