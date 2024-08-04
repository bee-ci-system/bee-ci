import Image from 'next/image';
import { staticAssets } from '../_utils/static-assets';

const GithubIcon = ({
  height,
  width,
  className,
}: {
  height?: number;
  width?: number;
  className?: string;
}) => (
  <Image
    src={staticAssets.icons.GITHUB}
    height={height || 18}
    width={width || 18}
    alt='github icon'
    className={className}
  />
);

export { GithubIcon };
