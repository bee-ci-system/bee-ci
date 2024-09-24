import { cookies } from 'next/headers';
import { NextResponse, type NextRequest } from 'next/server';
import { documentationRoutes, routes } from './app/_utils/routes';

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  if (pathname === '/docs') {
    const redirectUrl = new URL(documentationRoutes[0].href, request.url);
    console.log('redirectUrl', redirectUrl);
    return NextResponse.redirect(redirectUrl);
  }

  const token = cookies().get('jwt');

  if (!token && pathname !== '/') {
    return NextResponse.redirect(new URL('/', request.url));
  }

  if (token && pathname === '/') {
    return NextResponse.redirect(new URL(routes.DASHBOARD, request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/public|public|.*\\.svg$|docs/.*).*)'],
};
