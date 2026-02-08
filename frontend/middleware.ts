import { cookies } from 'next/headers';
import { NextResponse, type NextRequest } from 'next/server';
import { documentationRoutes, routes } from './app/_utils/routes';

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  if (pathname === routes.LOG_OUT) {
    const response = NextResponse.redirect(new URL('/', request.url));
    response.cookies.set('jwt', '', {
      path: '/',
      expires: new Date(0),
      domain: 'pacia.tech',
      sameSite: 'lax',
    });
    return response;
  }

  if (pathname === '/docs') {
    const redirectUrl = new URL(documentationRoutes[0].href, request.url);
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
