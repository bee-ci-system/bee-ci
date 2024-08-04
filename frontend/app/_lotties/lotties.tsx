'use client';

import Lottie, { LottieComponentProps } from 'lottie-react';
import automateWorkflowLottie from './automate-workflow.json';
import contenerizationLotie from './containerization.json';
import githubLottie from './github.json';
import statisticsLottie from './statistics.json';

const AutomateWorkflowLottie = (
  props: Omit<LottieComponentProps, 'animationData'>,
) => {
  return <Lottie animationData={automateWorkflowLottie} {...props} />;
};

const ContainerizationLottie = (
  props: Omit<LottieComponentProps, 'animationData'>,
) => {
  return <Lottie animationData={contenerizationLotie} {...props} />;
};

const GithubLottie = (props: Omit<LottieComponentProps, 'animationData'>) => {
  return <Lottie animationData={githubLottie} {...props} />;
};

const StatisticsLottie = (
  props: Omit<LottieComponentProps, 'animationData'>,
) => {
  return <Lottie animationData={statisticsLottie} {...props} />;
};

export {
  AutomateWorkflowLottie,
  ContainerizationLottie,
  GithubLottie,
  StatisticsLottie,
};
