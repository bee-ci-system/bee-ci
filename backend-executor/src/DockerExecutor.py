# executor class that runs the code in docker container and stream logs to influxdb
import docker
import tarfile
import os
import logging
import time

from InfluxDBHandler import InfluxDBHandler
from structures.InfluxDBCredentials import InfluxDBCredentials
from structures.BuildInfo import BuildInfo
from structures.BuildConfig import BuildConfig


class ExecutorFailure(Exception):
    pass


class ExecutorTimeout(Exception):
    pass


class DockerExecutor:
    def __init__(self, influxdb_credentials: InfluxDBCredentials):
        self.client = docker.from_env()
        self.influxdbHandler = InfluxDBHandler(influxdb_credentials)
        self.logger = logging.getLogger(__name__)
        self.logger.info("DockerExecutor initialized")

    def copy_to(self, src: str, container: docker.models.containers.Container):
        self.logger.debug("Copying file to container")
        srcname = os.path.basename(src)
        tar = tarfile.open(src + ".tar", mode="w")
        try:
            tar.add(srcname)
        except Exception as e:
            self.logger.error("Error while copying file to container")
            raise ExecutorFailure from e
        finally:
            tar.close()
        self.logger.debug("Tar file created %s.tar", src)

        data = open(src + ".tar", "rb").read()

        container.put_archive(os.path.dirname("/tmp/run.sh"), data)

    def pull_image(self, image: str):
        try:
            self.client.images.pull(image)
        except docker.errors.APIError as e:
            self.logger.error("Image not found: %s", image)
            raise ExecutorFailure from e
        self.logger.debug('Image: "%s" pulled', image)

    def run_container(self, build_config: BuildConfig, build_info: BuildInfo):
        script_path = "run.sh"
        try:
            with open(script_path, "r", encoding="utf-8"):
                pass
        except FileNotFoundError as e:
            self.logger.error("Fatal error - File not found: %s", script_path)
            raise ExecutorFailure from e

        self.pull_image(build_config.image)

        container = self.client.containers.create(
            build_config.image, ["/bin/sh", "/tmp/run.sh"], detach=True
        )
        self.logger.info("Container created: %s", container.name)

        self.copy_to(script_path, container)

        container.start()
        start_time = time.time()
        timeout = build_config.timeout

        try:
            for line in container.logs(stream=True):
                self.logger.debug(line.strip())
                self.influxdbHandler.log_to_influxdb(
                    build_info.build_id, str(line.strip())
                )

                # Check for timeout
                if time.time() - start_time > timeout:
                    self.logger.error("Container execution timed out")
                    container.stop()
                    raise ExecutorTimeout("Container execution timed out")

        except Exception as e:
            self.logger.error("Error during container execution: %s", str(e))
            container.stop()
            raise e
        finally:
            exit_status = container.wait()
            if exit_status["StatusCode"] != 0:
                self.logger.error(
                    "Container exited with status code %s", exit_status["StatusCode"]
                )
                raise ExecutorFailure(
                    f"Container exited with status code {exit_status['StatusCode']}"
                )
            else:
                self.logger.info("Container exited successfully with status code 0")
            container.remove()
            self.logger.info("Container removed")
