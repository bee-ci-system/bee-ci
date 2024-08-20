#executor class that runs the code in docker container and stream logs to influxdb
import docker
from datetime import datetime, timezone
import tarfile
import sys
import os

from InfluxDBHandler import InfluxDBHandler
from InfluxDBCredentials import InfluxDBCredentials

# Rest of the code...
influxDBCredentials = InfluxDBCredentials(
influxdbBucket = "home",
influxdbOrg = "beeci",
influxdbToken = "F7EGfgf3ckz7gzpFA3ZHrmHQpSl1-8nW_IwwmgF8BGNz4p9lX6Ew-mwcl_ryLfCL-JNrE6srwytikagNeLavyw==",
influxdbUrl = "http://localhost:8086"
)

# Rest of the code...
class Executor:
    def __init__(self):
        self.client = docker.from_env()
        self.influxdbHandler = InfluxDBHandler(influxDBCredentials)
    def copy_to(self, src, container):
        srcname = os.path.basename(src)
        tar = tarfile.open(src + '.tar', mode='w')
        try:
            tar.add(srcname)
        finally:
            tar.close()
        print("Tar file created")

        data = open(src + '.tar', 'rb').read()

        container.put_archive(os.path.dirname("/tmp/run.sh"), data)

    def create_container(self, script_path):
        try:
            with open(script_path, "r") as f:
                pass
        except FileNotFoundError:
            print("File not found")
            sys.exit(1)
        if not os.access(script_path, os.X_OK):
            print("File is not executable")
            sys.exit(1)

        container = self.client.containers.create("alpine",["/bin/sh", "/tmp/run.sh"], detach=True)
        print("Container created: "+ container.name)

        print("Copying file to container")
        self.copy_to(script_path, container)

        container.start()

        for line in container.logs(stream=True):
            #print(line.strip())
            self.influxdbHandler.log_to_influxdb(str(line.strip()))

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python Executor.py <script_path>")
        sys.exit(1)
    script_path = sys.argv[1]
    executor = Executor()
    executor.create_container(script_path)

    #printing logs from influx db
    fluxtable = executor.influxdbHandler.download_logs()
    for table in fluxtable:
        for record in table.records:
            print(record.values.get("_value"))