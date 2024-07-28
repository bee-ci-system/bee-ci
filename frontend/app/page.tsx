import { SideBySideLayout } from './_layouts/side-by-side-layout';
import {
  AutomateWorkflowLottie,
  ContainerizationLottie,
  GithubLottie,
  StatisticsLottie,
} from './_lotties/lotties';

export default function Home() {
  return (
    <SideBySideLayout>
      <div className='mx-3 flex flex-col gap-16 py-16 sm:mx-16'>
        <div className='flex w-full flex-col items-center rounded-md border-2 border-dashed border-beeci-yellow-100 bg-gradient-to-br from-gray-950 to-gray-800 py-4 pb-10 md:py-8 lg:flex-row'>
          <div className='w-1/2 sm:px-4 lg:w-2/5'>
            <ContainerizationLottie />
          </div>
          <div className='w-3/5 sm:px-4'>
            <p className='mb-6 text-center text-2xl font-medium lg:text-left'>
              🌐 Seamless GitHub Integration
            </p>
            <p className='text-center lg:text-left'>
              Automatically trigger CI actions with
              <strong className='ml-1'>every pull request</strong>. Stay in sync
              with your team and streamline code reviews and merges.
            </p>
          </div>
        </div>
        <div className='my-8 flex flex-col items-center px-4 lg:flex-row-reverse'>
          <div className='mb-4 w-1/2 px-4 pb-4 lg:w-2/5'>
            <AutomateWorkflowLottie loop={1} />
          </div>
          <div className='w-3/5'>
            <p className='mb-6 text-center text-2xl font-medium lg:text-left'>
              Automate Your Workflow 🛠️
            </p>
            <p className='text-center lg:text-left'>
              Define and manage your CI workflows easily with
              <strong className='ml-1'>YAML files</strong>. Customize every step
              of your build, test, and deploy processes.
            </p>
          </div>
        </div>
        <div className='flex flex-col items-center rounded-md border-2 border-dashed border-beeci-yellow-100 bg-gradient-to-br from-gray-950 to-gray-800 py-12 lg:flex-row'>
          <div className='w-4/5 px-16 pb-6 xs:w-5/12 sm:w-2/5 md:w-3/5 lg:w-2/5'>
            <GithubLottie />
          </div>
          <div className='mb-2 w-3/5 sm:px-4'>
            <p className='mb-6 text-center text-2xl font-medium lg:text-left'>
              Containerized Testing 🐳
            </p>
            <p className='text-center lg:text-left'>
              Run your tests in
              <strong className='ml-1'>isolated Docker containers</strong>.
              Ensure consistent and reliable results across different
              environments.
            </p>
          </div>
        </div>
        <div className='flex flex-col items-center py-4 lg:flex-row-reverse'>
          <div className='px-12 pb-8 sm:w-1/2 lg:w-2/5'>
            <StatisticsLottie loop={1} />
          </div>
          <div className='w-3/5 px-4'>
            <p className='mb-6 text-center text-2xl font-medium lg:text-left'>
              📊 Real-time Monitoring
            </p>
            <p className='text-center lg:text-left'>
              Keep track of your build and test results with
              <strong className='ml-1'>detailed logs and statistics</strong>.
              Monitor progress and quickly identify issues.
            </p>
          </div>
        </div>
      </div>
    </SideBySideLayout>
  );
}
