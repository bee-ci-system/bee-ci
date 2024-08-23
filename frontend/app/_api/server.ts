interface GetUserDto {
  name: string;
}

export function getUserServer(): GetUserDto {
  return { name: 'John Doe' };
}
