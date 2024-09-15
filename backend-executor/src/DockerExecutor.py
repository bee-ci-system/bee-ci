# executor class that runs the code in docker container and stream logs to influxdb
import docker
from datetime import datetime, timezone
import tarfile
import sys
import os
import logging

from InfluxDBHandler import InfluxDBHandler
from structures.InfluxDBCredentials import InfluxDBCredentials
from structures.BuildInfo import BuildInfo
from structures.BuildConfig import BuildConfig


influxDBCredentials = InfluxDBCredentials(
    influxdbBucket="home",
    influxdbOrg="beeci",
    influxdbToken="9uNp_AJQknsl8OWY65VGyAVZ0wpLXrm9Ep9_4L4-LJJWkP4HJxQvgMCd0vIElfFVU-9cIMdPgPGuUZvaDJsn5g==",
    influxdbUrl="http://localhost:8086",
)


class ExecutorFailure(Exception):
    pass


class DockerExecutor:
    def __init__(self):
        self.client = docker.from_env()
        self.influxdbHandler = InfluxDBHandler(influxDBCredentials)
        self.logger = logging.getLogger(__name__)

    def copy_to(self, src, container: docker.models.containers.Container):
        self.logger.debug("Copying file to container")
        srcname = os.path.basename(src)
        tar = tarfile.open(src + ".tar", mode="w")
        try:
            tar.add(srcname)
        except:
            self.logger.error("Error while copying file to container")
            raise ExecutorFailure
        finally:
            tar.close()
        self.logger.debug("Tar file created" + src + ".tar")

        data = open(src + ".tar", "rb").read()

        container.put_archive(os.path.dirname("/tmp/run.sh"), data)

    def pull_image(self, image):
        try:
            self.client.images.pull(image)
        except docker.errors.APIError:
            self.logger.error("Image not found: " + image)
            raise ExecutorFailure
        self.logger.debug('Image: "' + image + '" pulled')

    def run_container(self, build_config: BuildConfig, build_info: BuildInfo):
        script_path = "run.sh"
        try:
            with open(script_path, "r") as f:
                pass
        except FileNotFoundError:
            self.logger.error("Fatal error - File not found: " + script_path)
            raise ExecutorFailure

        self.pull_image(build_config.image)

        container = self.client.containers.create(
            build_config.image, ["/bin/sh", "/tmp/run.sh"], detach=True
        )
        self.logger.info("Container created: " + container.name)

        self.copy_to(script_path, container)

        container.start()

        for line in container.logs(stream=True):
            self.logger.debug(line.strip())
            self.influxdbHandler.log_to_influxdb(build_info.id, str(line.strip()))

        container.remove()
        self.logger.info("Container removed")
