# executor class that runs the code in docker container and stream logs to influxdb
import docker
from datetime import datetime, timezone
import tarfile
import sys
import os

from InfluxDBHandler import InfluxDBHandler
from structures.InfluxDBCredentials import InfluxDBCredentials
from structures.BuildInfo import BuildInfo
from structures.BuildConfig import BuildConfig

# Rest of the code...
influxDBCredentials = InfluxDBCredentials(
    influxdbBucket="home",
    influxdbOrg="beeci",
    influxdbToken="9uNp_AJQknsl8OWY65VGyAVZ0wpLXrm9Ep9_4L4-LJJWkP4HJxQvgMCd0vIElfFVU-9cIMdPgPGuUZvaDJsn5g==",
    influxdbUrl="http://localhost:8086",
)


# Rest of the code...
class DockerExecutor:
    def __init__(self):
        self.client = docker.from_env()
        self.influxdbHandler = InfluxDBHandler(influxDBCredentials)

    def copy_to(self, src, container: docker.models.containers.Container):
        srcname = os.path.basename(src)
        tar = tarfile.open(src + ".tar", mode="w")
        try:
            tar.add(srcname)
        finally:
            tar.close()
        print("Tar file created")

        data = open(src + ".tar", "rb").read()

        container.put_archive(os.path.dirname("/tmp/run.sh"), data)

    def run_container(self, build_config: BuildConfig, build_info: BuildInfo):
        script_path = "run.sh"
        try:
            with open(script_path, "r") as f:
                pass
        except FileNotFoundError:
            print("Fatal error - File not found")
            sys.exit(1) 

        self.client.images.pull(build_config.image)
        print("Image pulled")
        container = self.client.containers.create(
            build_config.image, ["/bin/sh", "/tmp/run.sh"], detach=True
        )
        print("Container created: " + container.name)

        print("Copying file to container")
        self.copy_to(script_path, container)

        container.start()

        for line in container.logs(stream=True):
            self.influxdbHandler.log_to_influxdb(build_info.id, str(line.strip()))

        container.remove()
        print("Container removed")
