from datetime import datetime, timezone
import influxdb_client
from influxdb_client.client.write_api import SYNCHRONOUS
from structures.InfluxDBCredentials import InfluxDBCredentials


class InfluxDBHandler:
    def __init__(self, cred):
        self.cred = cred
        self.client = influxdb_client.InfluxDBClient(
            url=self.cred.url, token=self.cred.token, org=self.cred.org
        )
        self.write_api = self.client.write_api(write_options=SYNCHRONOUS)

    def log_to_influxdb(self, build_id, message):
        p = (
            influxdb_client.Point(build_id)
            .time(time=datetime.now(tz=timezone.utc))
            .field("Log", message)
        )
        self.write_api.write(bucket=self.cred.bucket, org=self.cred.org, record=p)

    def download_logs(self, build_id):
        query = f'from(bucket: "{self.cred.bucket}") |> range(start: -1h) |> filter(fn: (r) => r["_measurement"] == "{build_id}")'
        tables = self.client.query_api().query(query, org=self.cred.org)
        return tables
