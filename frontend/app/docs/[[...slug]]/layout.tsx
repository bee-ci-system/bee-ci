import { Leftbar } from '../_components/leftbar';
import { Navbar } from '../_components/navbar';

export default function DocsLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className='w-full'>
      <Navbar />
      <div className='flex w-full items-start px-8'>
        <Leftbar />
        <div className='flex flex-grow justify-between'>{children}</div>
      </div>
    </div>
  );
}
