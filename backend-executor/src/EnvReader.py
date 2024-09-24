import os
import sys
import logging

logger = logging.getLogger(__name__)
class EnvReader:
    @staticmethod
    def get_env_variables():
        # Get environment variables
        db_host = os.getenv("DB_HOST")
        db_port = os.getenv("DB_PORT")
        db_user = os.getenv("DB_USER")
        db_password = os.getenv("DB_PASSWORD")
        db_name = os.getenv("DB_NAME")

        influxdb_url = os.getenv("INFLUXDB_URL")
        influxdb_org = os.getenv("INFLUXDB_ORG")
        influxdb_bucket = os.getenv("INFLUXDB_BUCKET")
        influxdb_token = os.getenv("INFLUXDB_TOKEN")

        missing_env_vars = []
        if not db_host:
            missing_env_vars.append("DB_HOST")
        if not db_port:
            missing_env_vars.append("DB_PORT")
        if not db_user:
            missing_env_vars.append("DB_USER")
        if not db_password:
            missing_env_vars.append("DB_PASSWORD")
        if not db_name:
            missing_env_vars.append("DB_NAME")
        if not influxdb_url:
            missing_env_vars.append("INFLUXDB_URL")
        if not influxdb_org:
            missing_env_vars.append("INFLUXDB_ORG")
        if not influxdb_bucket:
            missing_env_vars.append("INFLUXDB_BUCKET")
        if not influxdb_token:
            missing_env_vars.append("INFLUXDB_TOKEN")

        if missing_env_vars:
            logger.error(f"Missing environment variables: {', '.join(missing_env_vars)}")
            sys.exit(1)

        logger.info(f"DB Host: {db_host}:{db_port}, DB User: {db_user}, DB Name: {db_name}")
        logger.info(f"InfluxDB Url: {influxdb_url}, Org: {influxdb_org}, Bucket: {influxdb_bucket}")

        return {
            "db_host": db_host,
            "db_port": db_port,
            "db_user": db_user,
            "db_password": db_password,
            "db_name": db_name,
            "influxdb_url": influxdb_url,
            "influxdb_org": influxdb_org,
            "influxdb_bucket": influxdb_bucket,
            "influxdb_token": influxdb_token,
        }