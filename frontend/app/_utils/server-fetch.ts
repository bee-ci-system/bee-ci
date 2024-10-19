import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { apiBaseUrl } from '../_utils/constants';
import { routes } from './routes';

export const serverFetch = async (
  endpoint: string,
  options: RequestInit = {},
) => {
  const token = cookies().get('jwt');

  const headers: HeadersInit = new Headers({
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*',
    ...options.headers,
  });

  if (token) {
    headers.set('Cookie', `jwt=${token.value}`);
  }

  const url = `${apiBaseUrl}${endpoint}`;
  console.log(`serverFetch: will fetch from: ${url}`);
  const response = await fetch(url, {
    ...options,
    headers,
    cache: 'no-store',
  });

  console.log(
    `serverFetch: did fetch from: ${url}, status: ${response.status}`,
  );

  if (response.status === 401) {
    return redirect(routes.LOG_OUT);
  }

  return response;
};
