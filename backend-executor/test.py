import docker
from datetime import datetime, timezone
import tarfile
import sys
import os
import influxdb_client
from influxdb_client.client.write_api import SYNCHRONOUS

client = docker.from_env()
influxdbBucket = "home"
influxdbOrg = "beeci"
influxdbToken = "oDPScfvTrwxr3kKbb22GNpgQsRQ57idolaLxD3UBIARU9AbWmjz6cb8yDcJSUHrHrRQ-omeOhmh09Rn7RSEStg=="
influxdbUrl = "http://localhost:8086"
influxdbClient = influxdb_client.InfluxDBClient(
            url=influxdbUrl,
            token=influxdbToken,
            org=influxdbOrg
        )

def copy_to(src, container):

    os.chdir(os.path.dirname(src))
    srcname = os.path.basename(src)
    tar = tarfile.open(src + '.tar', mode='w')
    try:
        tar.add(srcname)
    finally:
        tar.close()

    data = open(src + '.tar', 'rb').read()

    container.put_archive(os.path.dirname("/tmp/run.sh"), data)

def write_to_influxdb(data):
    write_api = influxdbClient.write_api(write_options=SYNCHRONOUS)
    write_api.write(bucket=influxdbBucket, org=influxdbOrg, record=data)

#create main function
def main():
    #read arguments, there should be one with script path, use argument parser
    if len(sys.argv) != 2:
        print("Usage: python test.py <script_path>")
        sys.exit(1)
    #get script path
    script_path = sys.argv[1]
    try:
        with open(script_path, "r") as f:
            pass
    except FileNotFoundError:
        print("File not found")
        sys.exit(1)
    if not os.access(script_path, os.X_OK):
        print("File is not executable")
        sys.exit(1)

    container = client.containers.create("alpine",["/bin/sh", "/tmp/run.sh"], detach=True)
    print("Container created: "+ container.name)

    print("Copying file to container")
    copy_to(script_path, container)

    container.start()

    counter = 1
    for line in container.logs(stream=True):
        print(line.strip())
        p = influxdb_client.Point("Log").tag("Test", counter).time(time=datetime.now(tz=timezone.utc)).field("message", str(line.strip()))
        write_to_influxdb(p)
        counter += 1

    container.remove()
    print("Container removed")

if __name__ == "__main__":
    main()

