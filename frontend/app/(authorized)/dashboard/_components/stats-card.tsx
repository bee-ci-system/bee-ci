import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/app/_components/card';
import { ReactNode } from 'react';

const StatsCard = ({
  title,
  value,
  percent,
  icon,
}: {
  title: string;
  value: number;
  icon: ReactNode;
  percent?: number;
}) => {
  return (
    <Card className='border-beeci-yellow-600'>
      <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
        <CardTitle className='text-sm font-medium'>{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className='text-2xl font-bold'>{value}</div>
        {percent && (
          <p className='text-right text-xs text-muted-foreground'>{percent}%</p>
        )}
      </CardContent>
    </Card>
  );
};

export { StatsCard };
