import { cookies } from 'next/headers';
import { NextResponse, type NextRequest } from 'next/server';
import { documentationRoutes, routes } from './app/_utils/routes';

export function middleware(request: NextRequest) {
  console.log(
    `DEBUG middleware: started. request: ${console.log(JSON.stringify(request, null, 4))}`,
  );
  const { pathname } = request.nextUrl;

  if (pathname === '/docs') {
    console.log(`DEBUG middleware: entered pathname === /docs`);
    const redirectUrl = new URL(documentationRoutes[0].href, request.url);
    console.log('redirectUrl', redirectUrl);
    return NextResponse.redirect(redirectUrl);
  }

  const token = cookies().get('token');
  console.log(`DEBUG middleware: token cookie: ${token}`);
  console.log(`DEBUG middleware: jwt cookie: ${cookies().get('jwt')}`);

  if (!token && pathname !== '/') {
    console.log(`DEBUG middleware: entered !token && pathname !== '/'`);
    return NextResponse.redirect(new URL('/', request.url));
  }

  if (token && pathname === '/') {
    console.log(`DEBUG middleware: entered token && pathname === '/'`);
    return NextResponse.redirect(new URL(routes.DASHBOARD, request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/public|public|.*\\.svg$|docs/.*).*)'],
};
