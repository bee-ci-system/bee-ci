import { BookOpen } from 'lucide-react';
import Image from 'next/image';
import { ReactNode } from 'react';
import { Button } from '../_components/button';
import { GithubIcon } from '../_icons/github-icon';
import { staticAssets } from '../_utils/static-assets';

export const SideBySideLayout = ({ children }: { children: ReactNode }) => (
  <div className='h-screen'>
    <div className='flex h-full grow flex-col md:flex-row'>
      <div className='border-beeci-yellow-600 animated-background to-beeci-yellow-950 flex h-full flex-col justify-center gap-6 border-r-[1px] bg-gradient-to-br from-gray-950 from-10% via-gray-900 via-80% to-100% px-14 pb-4 pt-6 sm:pt-0 md:w-5/12'>
        <div>
          <Image
            src={staticAssets.logo.BEECI_LOGO_DARK_MODE}
            height={151}
            width={209}
            alt='beeci logo'
          />
        </div>
        <div>
          <p className='font-display flex flex-col whitespace-pre-wrap text-4xl/tight font-light text-white'>
            Open-source CI system
            <span className='text-beeci-yellow-400'>
              for modern developers{' '}
            </span>
          </p>
          <p className='mt-4 text-sm/6'>
            BeeCI is a streamlined CI tool you can rely on to automate your
            build, test, and deployment processes effortlessly. With simple YAML
            workflow definitions and seamless GitHub integration, it's
            efficient, user-friendly, and designed to make your development
            lifecycle buzz 🐝 with productivity.
          </p>
          <Button className='mt-6 w-1/2'>Begin here </Button>
        </div>
        <div className='flex w-full flex-col items-start text-sm md:flex-row md:gap-4 lg:gap-8'>
          <Button variant='link' className='px-0'>
            <GithubIcon className='mr-1 px-0' />
            Documentation
          </Button>
          <Button variant='link' className='px-0'>
            <BookOpen size={18} className='mr-1' /> Github
          </Button>
        </div>
      </div>
      <div className='md:w-7/12'>
        <div className='h-full p-4'>{children}</div>
      </div>
    </div>
  </div>
);
