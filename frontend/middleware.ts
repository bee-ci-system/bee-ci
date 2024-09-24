import { cookies } from 'next/headers';
import { NextResponse, type NextRequest } from 'next/server';
import { documentationRoutes, routes } from './app/_utils/routes';

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  console.log(`
-----
DEBUG: New request! Pathname: ${pathname}. Content:
${JSON.stringify(request, null, 4)}
-----
`);

  if (pathname === '/log-out') {
    // Delete token on logout
    console.log(`DEBUG: entered pathname === /docs`);
    const response = NextResponse.redirect(new URL('/', request.url));
    response.cookies.set('jwt', '', { path: '/', expires: new Date(0) });
    return response;
  }

  if (pathname === '/docs') {
    console.log(`DEBUG: entered pathname === /docs`);
    const redirectUrl = new URL(documentationRoutes[0].href, request.url);
    console.log('DEBUG: redirectUrl:', redirectUrl);
    return NextResponse.redirect(redirectUrl);
  }

  const token = cookies().get('jwt');
  console.log(`DEBUG: token: ${JSON.stringify(token)}`);

  if (!token && pathname !== '/') {
    // For non-authenticated users, always redirect to landing
    console.log(`DEBUG: entered !token && pathname !== '/'`);
    return NextResponse.redirect(new URL('/', request.url));
  }

  if (token && pathname === '/') {
    // For authenticated users, always redirect from landing to dashboard
    console.log(`DEBUG: entered token && pathname === '/'`);
    return NextResponse.redirect(new URL(routes.DASHBOARD, request.url));
  }

  console.log(`DEBUG: entered default return. END REQUEST -->`);
  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/public|public|.*\\.svg$|docs/.*).*)'],
};
