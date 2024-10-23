import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { routes } from './routes';

export const serverFetch = async (
  endpoint: string,
  options: RequestInit = {},
) => {
  let apiBaseUrl = process.env.API_URL_SERVER_OVERRIDE;
  if (!apiBaseUrl) {
    apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  }

  if (!apiBaseUrl) {
    throw new Error(
      'Neither NEXT_PUBLIC_API_BASE_URL nor API_URL_SERVER_OVERRIDE is set',
    );
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

  const url = `${apiBaseUrl}${endpoint}`;
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
