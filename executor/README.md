# Executor

## Overview

The Executor is a component of the BeeCI system responsible for running build jobs inside Docker containers and streaming logs to InfluxDB. It pulls build configurations from the database, executes the specified commands, and logs the output.

## Features

- **Docker Integration**: Runs build jobs inside Docker containers.
- **InfluxDB Logging**: Streams logs to InfluxDB for monitoring and analysis.
- **Error Handling**: Handles various execution errors and timeouts.

## Prerequisites

- Docker installed on the host machine.
- Python 3.12 or higher.
- [uv](https://github.com/astral-sh/uv) - Python package manager

## Docker run
Dockerfile is available for building and running the executor in a container. The container is built with the following command:
```sh
docker build -t bee-ci-executor .
```

The container can be run with the following command:
```sh
docker run -d --env-file .env bee-ci-executor
```

## Manual installation

1. Clone the repository:
    ```sh
    git clone https://github.com/bee-ci-system/bee-ci.git
    cd bee-ci/executor
    ```

2. Install the required Python packages with uv:
    ```sh
    uv sync
    ```

## Configuration

Ensure that the environment variables are set correctly. The executor reads environment variables, all needed variables are defined in .env.sample at main directory

InfluxDb2 and Postgres databases needs to be running on adresses defined in .env.sample.

## Usage

```sh
python3 src/main.py
```