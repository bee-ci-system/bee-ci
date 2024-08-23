import Image from 'next/image';
import { staticAssets } from '../_utils/static-assets';

const Logo = ({
  type = 'short',
  size = 18,
}: {
  type: 'short' | 'long';
  size: number;
}) => {
  const logoConfig = {
    short: {
      light: {
        src: staticAssets.logo.BEECI_LOGO_SHORT_LIGHT_MODE,
        className: 'dark:hidden',
      },
      dark: {
        src: staticAssets.logo.BEECI_LOGO_SHORT_DARK_MODE,
        className: 'hidden dark:block',
      },
    },
    long: {
      light: {
        src: staticAssets.logo.BEECI_LOGO_LIGHT_MODE,
        className: 'dark:hidden',
      },
      dark: {
        src: staticAssets.logo.BEECI_LOGO_DARK_MODE,
        className: 'hidden dark:block',
      },
    },
  };

  const { light, dark } = logoConfig[type];

  return (
    <>
      <Image
        src={light.src}
        alt='bee-ci logo'
        height={size}
        width={type === 'short' ? size : size * 1.38}
        className={light.className}
      />
      <Image
        src={dark.src}
        alt='bee-ci logo'
        height={size}
        width={type === 'short' ? size : size * 1.38}
        className={dark.className}
      />
    </>
  );
};

export { Logo };
