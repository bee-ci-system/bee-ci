import { documentationRoutes, routes } from '@/app/_utils/routes';
import DocsBreadcrumb from '@/app/docs/_components/docs-breadcrumb';
import Toc from '@/app/docs/_components/toc';
import { notFound } from 'next/navigation';
import { cache, PropsWithChildren } from 'react';
import Pagination from '../_components/pagination';
import { getMarkdownForSlug } from '../_utils/markdown';

type PageProps = {
  params: { slug: string[] };
};

const cachedGetMarkdownForSlug = cache(getMarkdownForSlug);

async function DocumentationPage({ params: { slug = [] } }: PageProps) {
  const pathName = slug.join('/');
  const res = await cachedGetMarkdownForSlug(pathName);

  if (!res) notFound();
  return (
    <div className='flex w-full flex-grow justify-between gap-8'>
      <div className='flex-grow pt-12'>
        <DocsBreadcrumb paths={slug} />
        <Markdown>
          <h1>{res.frontmatter.title}</h1>
          <p className='text-[16.5px] text-muted-foreground'>
            {res.frontmatter.description}
          </p>
          <div>{res.content}</div>
          <Pagination pathname={'/docs/' + pathName} />
        </Markdown>
      </div>
      <Toc path={pathName} />
    </div>
  );
}

function Markdown({ children }: PropsWithChildren) {
  return (
    <div className='prose-code:font-code prose prose-zinc w-fit pt-2 dark:prose-invert prose-headings:scroll-m-20 prose-code:rounded-md prose-code:bg-neutral-100 prose-code:p-1 prose-code:text-sm prose-code:leading-6 prose-code:text-neutral-800 prose-code:before:content-none prose-code:after:content-none prose-pre:border prose-pre:bg-neutral-100 dark:prose-code:bg-neutral-900 dark:prose-code:text-white dark:prose-pre:bg-neutral-900 sm:mx-auto'>
      {children}
    </div>
  );
}

export async function generateMetadata({ params: { slug = [] } }: PageProps) {
  const pathName = slug.join('/');
  const res = await cachedGetMarkdownForSlug(pathName);
  if (!res) return null;
  const { frontmatter } = res;
  return {
    title: frontmatter.title,
    description: frontmatter.description,
  };
}

export function generateStaticParams() {
  return documentationRoutes.map((item) => {
    const slug = item.href
      .replace(routes.DOCUMENTATION, '')
      .split('/')
      .filter((part) => part !== '');

    return { slug };
  });
}

export default DocumentationPage;
