import type { Metadata } from 'next';
import { Montserrat } from 'next/font/google';
import { ThemeProvider } from './_utils/theme-provider';
import './_styles/globals.css';
import ReactQueryProvider from './_utils/react-query-provider';

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
    <html lang='en' suppressHydrationWarning>
      <body className={inter.className}>
        <ReactQueryProvider>
          <ThemeProvider
            attribute='class'
            defaultTheme='system'
            enableSystem
            disableTransitionOnChange
          >
            {children}
          </ThemeProvider>
        </ReactQueryProvider>
      </body>
    </html>
  );
}
