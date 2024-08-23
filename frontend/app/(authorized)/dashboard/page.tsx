import { getUserServer } from '@/app/_api/server';
import { BadgeCheck, BadgeX, Computer } from 'lucide-react';
import { getDashboardData } from './_api/server';
import { PipelinesCard } from './_components/pipelines-card';
import { MyRepositoriesCard } from './_components/repositories-card';
import { StatsCard } from './_components/stats-card';
import { calculatePercent } from './_utils/calculate-percent';

const DashboardPage = async () => {
  const [userInfo, dashboardData] = await Promise.all([
    getUserServer(),
    getDashboardData(),
  ]);

  return (
    <div className='flex min-h-screen w-full flex-col'>
      <main className='flex flex-1 flex-col gap-4 p-4 md:gap-8 md:p-8'>
        <h1 className='prose-2xl ml-4 sm:-mt-16'>Hello {userInfo.name}!</h1>
        <div className='grid gap-4 md:grid-cols-2 md:gap-8 lg:grid-cols-3'>
          <StatsCard
            title='Total pipelines'
            value={dashboardData.stats.totalPipelines}
            icon={<Computer className='h-4 w-4 text-muted-foreground' />}
          />
          <StatsCard
            title='Successful pipelines'
            value={dashboardData.stats.successfulPipelines}
            icon={<BadgeCheck className='h-4 w-4 text-muted-foreground' />}
            percent={calculatePercent(
              dashboardData.stats.successfulPipelines,
              dashboardData.stats.totalPipelines,
            )}
          />
          <StatsCard
            title='Unsuccessful pipelines'
            value={dashboardData.stats.unsuccessfulPipelines}
            icon={<BadgeX className='h-4 w-4 text-muted-foreground' />}
            percent={calculatePercent(
              dashboardData.stats.unsuccessfulPipelines,
              dashboardData.stats.totalPipelines,
            )}
          />
        </div>
        <div className='grid gap-4 md:gap-8 lg:grid-cols-2 xl:grid-cols-3'>
          <MyRepositoriesCard repositories={dashboardData.repositories} />
          <PipelinesCard pipelines={dashboardData.pipelines} />
        </div>
      </main>
    </div>
  );
};

export default DashboardPage;
