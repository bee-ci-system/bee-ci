'use client';

import { FolderClosed, Home, Library, LogOut, PanelLeft } from 'lucide-react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '../_utils/cn';
import { routes } from '../_utils/routes';
import { Button } from './button';
import { Sheet, SheetContent, SheetTrigger } from './sheet';
import { ThemeToggler } from './theme-toggler';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from './tooltip';

const getNavLinkClassNames = (
  currentPathname: string,
  targetPathname: string,
) => {
  return cn(
    currentPathname === targetPathname
      ? 'bg-beeci-yellow-500 dark:bg-accent dark:text-beeci-yellow-500 text-white'
      : 'text-muted-foreground hover:text-white hover:bg-beeci-yellow-500',
    'flex h-9 w-9 items-center justify-center rounded-lg transition-colors md:h-8 md:w-8',
  );
};

const NavBar = () => {
  const currentPathname = usePathname();
  return (
    <>
      <aside className='fixed inset-y-0 left-0 z-10 hidden w-14 flex-col border-r bg-background sm:flex'>
        <TooltipProvider>
          <nav className='flex flex-col items-center gap-4 px-2 sm:py-5'>
            <Tooltip>
              <TooltipTrigger asChild>
                <Link
                  href={routes.DASHBOARD}
                  className={getNavLinkClassNames(
                    currentPathname,
                    routes.DASHBOARD,
                  )}
                >
                  <Home className='h-5 w-5' />
                  <span className='sr-only'>Dashboard</span>
                </Link>
              </TooltipTrigger>
              <TooltipContent side='right'>Dashboard</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger asChild>
                <Link
                  href={routes.MY_REPOSITORIES}
                  className={getNavLinkClassNames(
                    currentPathname,
                    routes.MY_REPOSITORIES,
                  )}
                >
                  <FolderClosed className='h-5 w-5' />
                  <span className='sr-only'>My repositories</span>
                </Link>
              </TooltipTrigger>
              <TooltipContent side='right'>My repositories</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger asChild>
                <Link
                  href={routes.DOCUMENTATION}
                  className={getNavLinkClassNames(
                    currentPathname,
                    routes.DOCUMENTATION,
                  )}
                >
                  <Library className='h-5 w-5' />
                  <span className='sr-only'>Documentation</span>
                </Link>
              </TooltipTrigger>
              <TooltipContent side='right'>Documentation</TooltipContent>
            </Tooltip>
          </nav>
          <nav className='mt-auto flex flex-col items-center gap-4 px-2 sm:py-5'>
            <ThemeToggler />
            <Tooltip>
              <TooltipTrigger asChild>
                <Link href={routes.LOGOUT} prefetch={false}>
                  <Button variant='outline' size='icon'>
                    <LogOut className='h-5 w-5' />
                    <span className='sr-only'>Log out</span>
                  </Button>
                </Link>
              </TooltipTrigger>
              <TooltipContent side='right'>Log out</TooltipContent>
            </Tooltip>
          </nav>
        </TooltipProvider>
      </aside>
      <div className='flex flex-col sm:gap-4 sm:py-4 sm:pl-14'>
        <header className='sticky top-0 z-30 flex h-14 items-center justify-between gap-4 border-b bg-background px-4 sm:static sm:h-auto sm:border-0 sm:bg-transparent sm:px-6'>
          <Sheet>
            <SheetTrigger asChild>
              <Button size='icon' variant='outline' className='sm:hidden'>
                <PanelLeft className='h-5 w-5' />
                <span className='sr-only'>Toggle Menu</span>
              </Button>
            </SheetTrigger>
            <SheetContent side='left' className='sm:max-w-xs'>
              <nav className='mt-4 grid gap-6 text-lg font-medium'>
                <Link
                  href={routes.DASHBOARD}
                  className={cn(
                    'flex items-center gap-4 px-2.5 text-muted-foreground',
                    routes.DASHBOARD === currentPathname &&
                      'text-beeci-yellow-500',
                  )}
                >
                  <Home className='h-5 w-5' />
                  Dashboard
                </Link>
                <Link
                  href={routes.MY_REPOSITORIES}
                  className={cn(
                    'flex items-center gap-4 px-2.5 text-muted-foreground',
                    routes.MY_REPOSITORIES === currentPathname &&
                      'text-beeci-yellow-500',
                  )}
                >
                  <FolderClosed className='h-5 w-5' />
                  My repositories
                </Link>
                <Link
                  href={routes.DOCUMENTATION}
                  className='flex items-center gap-4 px-2.5 text-muted-foreground'
                >
                  <Library className='h-5 w-5' />
                  Documentation
                </Link>
              </nav>
            </SheetContent>
          </Sheet>

          <div className='flex flex-row items-center gap-4 sm:hidden'>
            <ThemeToggler />
            <Link href={routes.LOGOUT} prefetch={false}>
              <Button variant='outline' size='icon'>
                <LogOut className='h-5 w-5' />
                <span className='sr-only'>Log out</span>
              </Button>
            </Link>
          </div>
        </header>
      </div>
    </>
  );
};

export { NavBar };
