import { Button } from '@/app/_components/button';
import { Command, CommandInput } from '@/app/_components/command';
import { cn } from '@/app/_utils/cn';
import { X } from 'lucide-react';

export const Search = ({
  searchValue,
  searchPlaceholder,
  handleSearchChange,
  className,
}: {
  searchValue: string;
  className?: string;
  searchPlaceholder?: string;
  //   eslint-disable-next-line no-unused-vars
  handleSearchChange: (value: string) => void;
}) => (
  <Command
    className={cn(
      'relative z-[1] order-3 col-span-2 flex h-auto w-full min-w-40 items-center justify-center rounded-none focus-within:ring-gray-600 sm:order-none sm:w-auto sm:min-w-[320px] sm:border-t-0 md:border-t-0 [&_[cmdk-input-wrapper]]:w-full [&_[cmdk-input-wrapper]]:border-b-0',
      className,
    )}
    label='search'
    aria-label='search'
  >
    <CommandInput
      key='input'
      placeholder={searchPlaceholder || 'Search Name'}
      value={searchValue}
      onValueChange={(value) => handleSearchChange(value)}
    />
    <Button
      variant='ghost'
      className='absolute right-0 top-1/2 h-11 w-10 -translate-y-1/2 px-0 hover:bg-transparent'
      aria-label='clear'
      onClick={() => handleSearchChange('')}
    >
      <X strokeWidth={0.6} />
    </Button>
  </Command>
);
