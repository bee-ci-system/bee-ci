# BeeCI 🐝

## Overview

BeeCI is a Continuous Integration (CI) system designed to automate building, testing, and deploying software using containerization technology.\
BeeCI simplifies CI processes, enhancing efficiency and productivity for developers.

## Key Features

### YAML Workflow Definition

Define test workflows with YAML files for flexible and efficient CI processes.

### GitHub Integration

Use GitHub webhooks to automatically trigger CI processes on new pull requests.

### CI Automation

Automatically run the CI process, including:

- Cloning the repository
- Building a Docker container
- Running tests in the container
- Analyzing test results

### Containerization

Execute tests in dedicated Docker containers as defined in the YAML workflow.

### Database

Track test logs and monitor repositories for comprehensive CI management.

## Tech Stack

<p align="left"> <a href="https://nextjs.org/" target="_blank" rel="noreferrer"> <img src="https://cdn.worldvectorlogo.com/logos/nextjs-2.svg" alt="nextjs" width="40" height="40"/> </a> <a href="https://reactjs.org/" target="_blank" rel="noreferrer"> <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/react/react-original-wordmark.svg" alt="react" width="40" height="40"/> </a> <a href="https://tailwindcss.com/" target="_blank" rel="noreferrer"> <img src="https://www.vectorlogo.zone/logos/tailwindcss/tailwindcss-icon.svg" alt="tailwind" width="40" height="40"/> </a> <a href="https://www.typescriptlang.org/" target="_blank" rel="noreferrer"> <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/typescript/typescript-original.svg" alt="typescript" width="40" height="40"/> </a> </p>

## Run Locally

Clone the project

```bash
  git clone https://github.com/kacaleksandra/bee-ci
```

Go to the project directory

```bash
  cd bee-ci
```

Install dependencies

```bash
  pnpm install
```

Start the server

```bash
  pnpm dev
```

## Demo

You can see the application on Vercel:\
https://bee-ci.vercel.app
