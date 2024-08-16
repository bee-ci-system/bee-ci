import { NextResponse, type NextRequest } from 'next/server';
import { documentationRoutes } from './app/_utils/routes';

export function middleware(request: NextRequest) {
  if (request.nextUrl.pathname === '/docs') {
    const redirectUrl = new URL(documentationRoutes[0].href, request.url);
    return NextResponse.redirect(redirectUrl);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/docs'],
};
