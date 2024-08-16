import { BuildTile, BuildTileProps } from './_components/build-tile';

const mockBuilds: BuildTileProps[] = [
  {
    id: '1',
    status: 'queued',
    repoOwner: 'bartekpacia',
    repoName: 'dumb',
    commitName: 'feat: add new feature',
    commitSha: '76de3f7',
  },
  {
    id: '2',
    status: 'in_progress',
    repoOwner: 'kacaleksandra',
    repoName: 'sample',
    commitName: 'fix: broken tests',
    commitSha: '76de3f7',
  },
  {
    id: '3',
    status: 'success',
    repoOwner: 'P3T3R',
    repoName: 'test',
    commitName: 'chore: update dependencies',
    commitSha: '76de3f7',
  },
];

const DashboardPage = () => {
  return (
    <main className='px-4 py-4'>
      <h1 className='pb-4 text-4xl'>Dashboard</h1>
      <main className='container flex flex-col gap-4'>
        {mockBuilds.map((build) => (
          <BuildTile
            key={build.id}
            id={build.id}
            status={build.status}
            repoOwner={build.repoOwner}
            repoName={build.repoName}
            commitName={build.commitName}
            commitSha={build.commitSha}
          />
        ))}
      </main>
    </main>
  );
};

export default DashboardPage;
