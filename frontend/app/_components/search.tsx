'use client';

import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTrigger,
} from '@/app/_components/dialog';
import { Input } from '@/app/_components/input';
import { ScrollArea } from '@/app/_components/scroll-area';
import { DialogTitle } from '@radix-ui/react-dialog';
import { FileTextIcon, SearchIcon } from 'lucide-react';
import { useMemo, useState } from 'react';
import { documentationRoutes } from '../_utils/routes';
import Anchor from './anchor';

export default function Search() {
  const [searchedInput, setSearchedInput] = useState('');
  const [isOpen, setIsOpen] = useState(false);

  const filteredResults = useMemo(
    () =>
      documentationRoutes.filter((item) =>
        item.title.toLowerCase().includes(searchedInput.toLowerCase()),
      ),
    [searchedInput],
  );

  return (
    <div>
      <Dialog
        open={isOpen}
        onOpenChange={(open) => {
          if (!open) setSearchedInput('');
          setIsOpen(open);
        }}
      >
        <DialogTrigger asChild>
          <div className='relative w-48 flex-1 cursor-pointer xs:w-64'>
            <SearchIcon className='absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-neutral-500 dark:text-neutral-400' />
            <Input
              className='h-9 w-full rounded-md border bg-muted pl-10 pr-4 text-sm shadow-sm'
              placeholder='Search documentation...'
              type='search'
            />
          </div>
        </DialogTrigger>
        <DialogContent className='top-[45%] max-w-[650px] p-0 sm:top-[38%]'>
          <DialogDescription className='sr-only'>
            Search in docs
          </DialogDescription>
          <DialogTitle className='sr-only'>Search</DialogTitle>
          <DialogHeader>
            <input
              value={searchedInput}
              onChange={(e) => setSearchedInput(e.target.value)}
              placeholder='Type something to search...'
              autoFocus
              className='h-14 border-b bg-transparent px-4 text-[15px] outline-none'
            />
          </DialogHeader>
          {filteredResults.length == 0 && searchedInput && (
            <p className='mx-auto mt-2 text-sm text-muted-foreground'>
              No results found for
              <span className='text-primary'>{`"${searchedInput}"`}</span>
            </p>
          )}
          <ScrollArea className='max-h-[350px]'>
            <div className='flex flex-col items-start gap-0.5 overflow-y-auto px-1 pb-4 sm:px-3'>
              {filteredResults.map((item) => (
                <DialogClose key={item.href} asChild>
                  <Anchor
                    className='flex w-full items-center gap-2.5 rounded-sm p-2.5 px-3 text-[15px] hover:bg-neutral-100 dark:hover:bg-neutral-900'
                    href={`${item.href}`}
                    activeClassName='dark:bg-neutral-900 bg-neutral-100'
                  >
                    <FileTextIcon className='h-[1.1rem] w-[1.1rem]' />
                    {item.title}
                  </Anchor>
                </DialogClose>
              ))}
            </div>
          </ScrollArea>
        </DialogContent>
      </Dialog>
    </div>
  );
}
