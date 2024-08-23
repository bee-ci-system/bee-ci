import { BookOpen } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import { ReactNode } from 'react';
import { Button, buttonVariants } from '../_components/button';
import { GithubIcon } from '../_icons/github-icon';
import { cn } from '../_utils/cn';
import { authUrl } from '../_utils/constants';
import { documentationRoutes } from '../_utils/routes';
import { staticAssets } from '../_utils/static-assets';

export const SideBySideLayout = ({ children }: { children: ReactNode }) => {
  return (
    <div className='flex h-screen grow flex-col md:flex-row'>
      <div className='animated-background flex h-full flex-col justify-center gap-6 border-b-[1px] border-beeci-yellow-600 bg-gradient-to-br from-gray-950 from-10% via-gray-900 via-80% to-beeci-yellow-950 to-100% px-14 pb-10 pt-16 md:w-5/12 md:border-b-0 md:border-r-[1px]'>
        <div>
          <Image
            src={staticAssets.logo.BEECI_LOGO_DARK_MODE}
            height={151}
            width={209}
            alt='bee-ci logo'
          />
        </div>
        <div>
          <p className='font-display flex flex-col whitespace-pre-wrap text-4xl/tight font-light text-white dark:text-primary'>
            Open-source CI system
            <span className='text-beeci-yellow-400'>for modern developers</span>
          </p>
          <p className='mt-4 text-sm/6 text-white dark:text-primary'>
            BeeCI is a streamlined CI tool you can rely on to automate your
            build, test, and deployment processes effortlessly. With simple YAML
            workflow definitions and seamless GitHub integration, it's
            efficient, user-friendly, and designed to make your development
            lifecycle buzz 🐝 with productivity.
          </p>
          <Button className='mt-6 w-1/2 min-w-max bg-white py-4 hover:bg-white dark:bg-white'>
            <a
              href={authUrl}
              className='text-black dark:text-primary-foreground'
            >
              Sign in with GitHub
            </a>
          </Button>
        </div>
        <div className='mb-4 flex w-full flex-col items-start text-sm lg:flex-row lg:gap-8'>
          <Link
            className={cn(
              buttonVariants({ variant: 'link' }),
              'px-0 text-white dark:text-primary',
            )}
            href={documentationRoutes[0].href}
          >
            <BookOpen size={18} className='mr-1 text-white dark:text-primary' />
            Read documentation
          </Link>
          <a
            href='https://github.com/kacaleksandra/bee-ci'
            target='_blank'
            className={cn(
              buttonVariants({ variant: 'link' }),
              'px-0 text-white dark:text-primary',
            )}
          >
            <GithubIcon className='mr-1 px-0' />
            View on GitHub
          </a>
        </div>
      </div>
      <div className='bg-gray-950 md:w-7/12 md:overflow-y-scroll'>
        <div className='h-full p-4'>{children}</div>
      </div>
    </div>
  );
};
