import type { Metadata } from 'next';
import { Montserrat } from 'next/font/google';
import { ThemeProvider } from './_utils/theme-provider';
import './_styles/globals.css';

const inter = Montserrat({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'BeeCI',
  description: 'BeeCI is a CI/CD platform for GitHub repositories.',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang='en'>
      <body className={inter.className}>
        <ThemeProvider
          attribute='class'
          defaultTheme='system'
          enableSystem
          disableTransitionOnChange
        >
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
