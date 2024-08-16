'use client';

import { Button } from '@/app/_components/button';
import { CheckIcon, CopyIcon } from 'lucide-react';
import { useState } from 'react';

export default function Copy({ content }: { content: string }) {
  const [isCopied, setIsCopied] = useState(false);

  async function handleCopy() {
    await navigator.clipboard.writeText(content);
    setIsCopied(true);

    setTimeout(() => {
      setIsCopied(false);
    }, 2000);
  }

  return (
    <Button variant='secondary' size='sm' onClick={handleCopy}>
      {isCopied ? (
        <CheckIcon className='h-3 w-3' />
      ) : (
        <CopyIcon className='h-3 w-3' />
      )}
    </Button>
  );
}
