import sys
import logging
import time
from DockerExecutor import DockerExecutor, ExecutorFailure, ExecutorTimeout
from DbPuller import DbPuller
from BuildConfigPuller import BuildConfigPuller
from BuildConfigAnalyzer import BuildConfigAnalyzer
from structures.BuildInfo import BuildInfo, BuildConclusion
from structures.InfluxDBCredentials import InfluxDBCredentials
from EnvReader import EnvReader

sleep_time = 10
logger = logging.getLogger("MainExecutor")


def print_logs(executor: DockerExecutor, build_id: int):
    # printing logs from influx db
    fluxtable = executor.influxdbHandler.download_logs(build_id)
    for table in fluxtable:
        for record in table.records:
            logger.info(record.values.get("_value"))

if __name__ == "__main__":
    if len(sys.argv) != 1:
        print("Usage: python main.py")
        sys.exit(1)
    logging.basicConfig(
        format="%(asctime)s %(name)s/%(levelname)s: %(message)s",
        datefmt="%H:%M:%S",
        level=logging.INFO,
    )
    env_vars = EnvReader.get_env_variables()

    db_puller = DbPuller(
        env_vars["db_host"],
        env_vars["db_port"],
        env_vars["db_name"],
        env_vars["db_user"],
        env_vars["db_password"],
    )
    influxdb_credentials = InfluxDBCredentials(
        env_vars["influxdb_bucket"],
        env_vars["influxdb_org"],
        env_vars["influxdb_token"],
        env_vars["influxdb_url"],
    )
    docker_executor = DockerExecutor(influxdb_credentials)
    while True:
        build_info = db_puller.pull_from_db()
        if not build_info:
            logger.info(
                "No available requests found in the database - sleeping for %d seconds",
                sleep_time,
            )
            time.sleep(sleep_time)
            continue

        config_data = BuildConfigPuller.pull_config(build_info)
        if not config_data:
            db_puller.update_conclusion(build_info.build_id, BuildConclusion.FAILURE)
            continue
        build_config = BuildConfigAnalyzer.analyze(config_data)
        if not build_config:
            db_puller.update_conclusion(build_info.build_id, BuildConclusion.FAILURE)
            continue

        try:
            docker_executor.run_container(build_config, build_info)
        except ExecutorFailure:
            logger.error("Failed to execute the build")
            db_puller.update_conclusion(build_info.build_id, BuildConclusion.FAILURE)
            continue
        except ExecutorTimeout:
            logger.error("Build execution timed out")
            db_puller.update_conclusion(build_info.build_id, BuildConclusion.TIMED_OUT)
            continue

        db_puller.update_conclusion(build_info.build_id, BuildConclusion.SUCCESS)
        print_logs(docker_executor, build_info.build_id)
