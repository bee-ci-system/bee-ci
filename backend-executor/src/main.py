import sys
import time
from DockerExecutor import DockerExecutor
from DbPuller import DbPuller
from BuildConfigPuller import BuildConfigPuller
from BuildConfigAnalyzer import BuildConfigAnalyzer
from structures.BuildInfo import BuildInfo

sleep_time = 10


def print_logs(executor):
    # printing logs from influx db
    fluxtable = executor.influxdbHandler.download_logs()
    for table in fluxtable:
        for record in table.records:
            print(record.values.get("_value"))

build_test = BuildInfo(
    id=1,
    repo_id=1,
    commit_sha="15f6c21b6fbba70d14ec14336cf79ce57fbc7ac2",
    commit_message="Update and rename .bee-ci.sh to .bee-ci.json ",
    status="queued",
    conclusion=None,
    created_at="2022-01-01T00:00:00Z",
    updated_at="2022-01-01T00:00:00Z",
    owner_name="bee-ci-system",
    repo_name="example-using-beeci",
)

if __name__ == "__main__":
    if len(sys.argv) != 1:
        print("Usage: python main.py")
        sys.exit(1)
    #script_path = sys.argv[1]
    db_puller = DbPuller()
    docker_executor = DockerExecutor()
    while True: 
        build_info = db_puller.pull_from_db()
        if not build_info:
            print(
                "No available requests found in the database - sleeping for "
                + str(sleep_time)
                + " second"
            )
            time.sleep(sleep_time)
            continue
        
        config_data = BuildConfigPuller.pull_config(build_test)
        if not config_data:
            db_puller.update_conclusion(build_info.id, "failed")
            continue
        build_config = BuildConfigAnalyzer.analyze(config_data)
        if not build_config:
            db_puller.update_conclusion(build_info.id, "failed")
            continue

        docker_executor.run_container(build_config, build_info)
        db_puller.update_conclusion(build_info.id, "success")
        print_logs(docker_executor)
