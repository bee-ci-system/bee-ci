import { RepositoriesTable } from './_components/repositories-table';

const MyRepositoriesPage = async () => {
  return (
    <div className='mx-auto mb-4 mt-4 flex min-h-[90vh] w-11/12'>
      <RepositoriesTable />
    </div>
  );
};

export default MyRepositoriesPage;
