import { ReactNode } from 'react';
import { NavBar } from '../_components/nav-bar';
import { ThemeToggler } from '../_components/theme-toggler';

const AuthorizedLayout = ({ children }: { children: ReactNode }) => {
  return (
    <div className='flex min-h-screen w-full flex-col bg-muted/40'>
      <NavBar />
      <div className='ml-auto mt-0 flex w-full flex-col pt-0 sm:gap-4 sm:pl-14'>
        <div className='hidden w-full pr-4 sm:flex sm:justify-end'>
          <ThemeToggler />
        </div>
        {children}
      </div>
    </div>
  );
};

export default AuthorizedLayout;
