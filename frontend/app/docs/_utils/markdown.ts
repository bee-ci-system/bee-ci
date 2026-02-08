import { promises as fs } from 'fs';
import path from 'path';
import Note from '@/app/_components/note';
import { Stepper, StepperItem } from '@/app/_components/stepper';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/app/_components/tabs';
import { documentationRoutes } from '@/app/_utils/routes';
import { compileMDX } from 'next-mdx-remote/rsc';
import rehypeAutolinkHeadings from 'rehype-autolink-headings';
import rehypeCodeTitles from 'rehype-code-titles';
import rehypePrism from 'rehype-prism-plus';
import rehypeSlug from 'rehype-slug';
import remarkGfm from 'remark-gfm';
import { visit } from 'unist-util-visit';
import Pre from '../_components/pre';

type MdxFrontmatter = {
  title: string;
  description: string;
};

const components = {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  pre: Pre,
  Note,
  Stepper,
  StepperItem,
};

export async function getMarkdownForSlug(slug: string) {
  try {
    const contentPath = getContentPath(slug);
    const rawMdx = await fs.readFile(contentPath, 'utf-8');
    return await compileMDX<MdxFrontmatter>({
      source: rawMdx,
      options: {
        parseFrontmatter: true,
        mdxOptions: {
          rehypePlugins: [
            preProcess,
            rehypeCodeTitles,
            rehypePrism,
            rehypeSlug,
            rehypeAutolinkHeadings,
            postProcess,
          ],
          remarkPlugins: [remarkGfm],
        },
      },
      components,
    });
  } catch (err) {
    console.log(err);
  }
}

export async function getTocs(slug: string) {
  const contentPath = getContentPath(slug);
  const rawMdx = await fs.readFile(contentPath, 'utf-8');

  const headingsRegex = /^(#{2,4})\s(.+)$/gm;
  let match;
  const extractedHeadings = [];
  while ((match = headingsRegex.exec(rawMdx)) !== null) {
    const headingLevel = match[1].length;
    const headingText = match[2].trim();
    const slug = sluggify(headingText);
    extractedHeadings.push({
      level: headingLevel,
      text: headingText,
      href: `#${slug}`,
    });
  }
  return extractedHeadings;
}

export function getPreviousNext(path: string) {
  const index = documentationRoutes.findIndex(({ href }) => href == path);
  return {
    prev: documentationRoutes[index - 1],
    next: documentationRoutes[index + 1],
  };
}

function sluggify(text: string) {
  const slug = text.toLowerCase().replace(/\s+/g, '-');
  return slug.replace(/[^a-z0-9-]/g, '');
}

function getContentPath(slug: string) {
  return path.join(process.cwd(), 'app/docs/_contents/', `${slug}.mdx`);
}

// for copying code blocks
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const preProcess = () => (tree: any) => {
  visit(tree, (node) => {
    if (node?.type === 'element' && node?.tagName === 'pre') {
      const [codeEl] = node.children;
      if (codeEl.tagName !== 'code') return;
      node.raw = codeEl.children?.[0].value;
    }
  });
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const postProcess = () => (tree: any) => {
  visit(tree, 'element', (node) => {
    if (node?.type === 'element' && node?.tagName === 'pre') {
      node.properties['raw'] = node.raw;
    }
  });
};
