from datetime import datetime, timezone
import influxdb_client
from influxdb_client.client.write_api import SYNCHRONOUS
from InfluxDBCredentials import InfluxDBCredentials

class InfluxDBHandler:
    def __init__(self, cred):
        self.cred = cred
        self.client = influxdb_client.InfluxDBClient(
            url=self.cred.url,
            token=self.cred.token,
            org=self.cred.org
        )
        self.write_api = self.client.write_api(write_options=SYNCHRONOUS)

    def log_to_influxdb(self, message):
        p = influxdb_client.Point("Log").time(time=datetime.now(tz=timezone.utc)).field("message", message)
        self.write_api.write(bucket=self.cred.bucket, org=self.cred.org, record=p)
        
    def download_logs(self):
        query = f'from(bucket: "{self.cred.bucket}") |> range(start: -1h)'
        tables = self.client.query_api().query(query, org=self.cred.org)
        return tables

