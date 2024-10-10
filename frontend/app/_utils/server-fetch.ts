import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { routes } from './routes';

export const serverFetch = async (
  endpoint: string,
  options: RequestInit = {},
) => {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;

  if (!baseUrl) {
    throw new Error('Base url is not established');
  }
  const token = cookies().get('jwt');

  const headers: HeadersInit = new Headers({
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*',
    ...options.headers,
  });

  if (token) {
    headers.set('Cookie', `jwt=${token.value}`);
  }

  const response = await fetch(`${baseUrl}${endpoint}`, {
    ...options,
    headers,
    cache: 'no-store',
  });

  if (response.status === 401) {
    return redirect(routes.LOG_OUT);
  }

  return response;
};
