import { Logo } from '@/app/_components/logo';
import { ThemeToggler } from '@/app/_components/theme-toggler';
import Link from 'next/link';
import Search from '../../_components/search';
import { SheetLeftbar } from './leftbar';

export function Navbar() {
  return (
    <nav className='sticky top-0 z-50 h-16 w-full border-b bg-opacity-5 px-2 backdrop-blur-xl backdrop-filter lg:px-4'>
      <div className='mx-auto flex h-full max-w-[1530px] items-center justify-between gap-2 p-2 sm:p-3'>
        <div className='flex items-center gap-5'>
          <SheetLeftbar />
          <div className='flex items-center gap-8'>
            <div className='hidden sm:flex'>
              <Link href='/' className='flex items-center gap-2.5'>
                <Logo type='short' size={18} />
                <h2 className='text-md font-bold'>bee-ci/docs</h2>
              </Link>
            </div>
          </div>
        </div>

        <div className='flex items-center gap-4'>
          <Search />
          <ThemeToggler />
        </div>
      </div>
    </nav>
  );
}
