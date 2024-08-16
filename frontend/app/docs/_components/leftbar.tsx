import { Button } from '@/app/_components/button';
import { Logo } from '@/app/_components/logo';
import { ScrollArea } from '@/app/_components/scroll-area';
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/app/_components/sheet';
import { routes } from '@/app/_utils/routes';
import { AlignLeftIcon } from 'lucide-react';
import Link from 'next/link';
import Anchor from '../../_components/anchor';

export function Leftbar() {
  return (
    <aside className='sticky top-16 hidden h-[92.75vh] min-w-[230px] flex-[0.9] flex-col overflow-y-auto lg:flex'>
      <ScrollArea className='py-4'>
        <Menu />
      </ScrollArea>
    </aside>
  );
}

export function SheetLeftbar() {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button variant='ghost' size='icon' className='flex lg:hidden'>
          <AlignLeftIcon />
        </Button>
      </SheetTrigger>
      <SheetContent className='flex flex-col gap-4 px-0' side='left'>
        <SheetTitle className='sr-only'>Menu</SheetTitle>
        <SheetDescription className='sr-only'>Menu</SheetDescription>
        <SheetHeader>
          <SheetClose className='px-5' asChild>
            <Link href='/' className='flex items-center gap-2.5'>
              <Logo type='short' size={18} />
              <h2 className='text-md font-bold'>bee-ci/docs</h2>
            </Link>
          </SheetClose>
        </SheetHeader>
        <ScrollArea className='flex flex-col gap-4'>
          <div className='mx-2 mt-3 flex flex-col gap-2 px-5'></div>
          <div className='mx-2 px-5'>
            <Menu isSheet />
          </div>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );
}

const MenuLink = ({ href, title }: { href: string; title: string }) => (
  <Anchor activeClassName='font-medium text-primary' href={`/docs/${href}`}>
    {title}
  </Anchor>
);

function Menu({ isSheet = false }) {
  return (
    <>
      {routes.documentationRoutes.map(({ href, items, title }) => {
        return (
          <div className='mt-5 flex flex-col gap-3' key={href}>
            <h4 className='font-medium sm:text-sm'>{title}</h4>
            <div className='ml-0.5 flex flex-col gap-3 text-neutral-800 dark:text-neutral-300/85 sm:text-sm'>
              {items.map((subItem) => {
                return isSheet ? (
                  <SheetClose key={`sheet-${href}${subItem.href}`} asChild>
                    <MenuLink
                      key={`${href}${subItem.href}`}
                      title={subItem.title}
                      href={`${href}${subItem.href}`}
                    />
                  </SheetClose>
                ) : (
                  <MenuLink
                    key={`${href}${subItem.href}`}
                    title={subItem.title}
                    href={`${href}${subItem.href}`}
                  />
                );
              })}
            </div>
          </div>
        );
      })}
    </>
  );
}
