import sys
import time
from Executor import Executor
from DbPuller import DbPuller
from BuildScriptPuller import BuildScriptPuller

sleep_time = 1


def print_logs(executor):
    # printing logs from influx db
    fluxtable = executor.influxdbHandler.download_logs()
    for table in fluxtable:
        for record in table.records:
            print(record.values.get("_value"))


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python main.py <script_path>")
        sys.exit(1)
    script_path = sys.argv[1]
    dbPuller = DbPuller()
    executor = Executor()
    while True:
        build = dbPuller.pull_from_db()
        
        if not build:
            print(
                "No available requests found in the database - sleeping for "
                + str(sleep_time)
                + " second"
            )
            time.sleep(sleep_time)
            continue
        executor.run_container(script_path, build)
        print_logs(executor)
