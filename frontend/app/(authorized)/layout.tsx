import { ReactNode } from 'react';
import { NavBar } from '../_components/nav-bar';

const AuthorizedLayout = ({ children }: { children: ReactNode }) => {
  return (
    <div className='flex min-h-screen w-full flex-col bg-muted/40'>
      <NavBar />
      <div className='ml-auto mt-0 flex w-full flex-col pt-0 sm:gap-4 sm:pl-14'>
        {children}
      </div>
    </div>
  );
};

export default AuthorizedLayout;

// Workaround for https://stackoverflow.com/q/79029275/7009800
export const dynamic = 'force-dynamic';
