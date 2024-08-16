'use client';

type BuildStatus = 'queued' | 'in_progress' | 'success' | 'failure';

interface BuildTileProps {
  id: string;
  status: BuildStatus;
  repoOwner: string;
  repoName: string;
  commitName: string;
  commitSha: string;
}

const BuildTile = (props: BuildTileProps) => {
  const { status, repoOwner, repoName, commitName, commitSha } = props;

  let color = '';

  if (status === 'success') {
    color = 'bg-green-500';
  } else if (status === 'in_progress') {
    color = 'bg-yellow-500';
  } else {
    color = 'bg-red-500';
  }

  return (
    <div className='flex rounded-lg bg-accent px-2 py-1'>
      <div className={`mx-4 my-auto rounded-full ${color} p-4`}></div>
      <div className='flex flex-col'>
        <div>
          <span className='text-xl'>{repoOwner}</span>
          <span className='mx-1 text-xl'>/</span>
          <span className='text-xl'>{repoName}</span>
          <br />
          <span>{commitSha}</span> • <span>{commitName}</span>
        </div>
      </div>
    </div>
  );
};

export type { BuildTileProps };
export { BuildTile };
