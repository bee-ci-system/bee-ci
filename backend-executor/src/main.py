import sys
import logging
import time
from DockerExecutor import DockerExecutor, ExecutorFailure
from DbPuller import DbPuller
from BuildConfigPuller import BuildConfigPuller
from BuildConfigAnalyzer import BuildConfigAnalyzer
from structures.BuildInfo import BuildInfo

sleep_time = 10
logger = logging.getLogger("MainExecutor")


def print_logs(executor, build_id):
    # printing logs from influx db
    fluxtable = executor.influxdbHandler.download_logs(build_id)
    for table in fluxtable:
        for record in table.records:
            logger.info(record.values.get("_value"))


build_test = BuildInfo(
    id=1,
    repo_id=1,
    commit_sha="0262a10fb0590f29471feed5ecf53b418b5b0d67",
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
    logging.basicConfig(
        format="%(asctime)s %(name)s/%(levelname)s: %(message)s",datefmt='%H:%M:%S', level=logging.INFO
    )
    db_puller = DbPuller()
    docker_executor = DockerExecutor()
    while True:
        build_info = db_puller.pull_from_db()
        if not build_info:
            logger.info(
                "No available requests found in the database - sleeping for "
                + str(sleep_time)
                + " seconds"
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

        try:
            docker_executor.run_container(build_config, build_info)
        except ExecutorFailure:
            logger.error("Failed to execute the build")
            db_puller.update_conclusion(build_info.id, "failed")
            continue
        except Exception as e:
            logger.fatal("An unexpected error occurred: " + str(e))
            db_puller.update_conclusion(build_info.id, "failed")
            continue

        db_puller.update_conclusion(build_info.id, "success")
        print_logs(docker_executor, build_info.id)
