import Image from 'next/image';
import { staticAssets } from '../_utils/static-assets';

export const Loader = () => (
  <div className='absolute inset-0 z-10 flex min-h-[100px] items-center justify-center bg-transparent'>
    <Image
      width={100}
      height={100}
      alt='loading'
      priority
      src={staticAssets.LOADER}
      unoptimized
    />
  </div>
);
