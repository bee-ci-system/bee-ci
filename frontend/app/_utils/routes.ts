export const routes = {
  LANDING: '/',
  DASHBOARD: '/dashboard',
  MY_REPOSITORIES: '/my-repositories',
  DOCUMENTATION: '/docs',
  LOG_OUT: '/logout',

  documentationRoutes: [
    {
      title: 'Getting Started',
      href: 'getting-started',
      items: [
        { title: 'Introduction', href: '/introduction' },
        { title: 'Installation', href: '/installation' },
      ],
    },
  ],
};

export const routeGenerators = {
  repository: (id: string) => `/repository/${id}`,
  pipeline: (id: string) => `/pipeline/${id}`,
};

export const documentationRoutes = routes.documentationRoutes
  .map(({ href, items }) =>
    items.map((link) => ({
      title: link.title,
      href: `${routes.DOCUMENTATION}/${href}${link.href}`,
    })),
  )
  .flat();
