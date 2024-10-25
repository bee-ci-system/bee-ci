class InfluxDBCredentials:
    def __init__(
        self,
        influxdbBucket: str,
        influxdbOrg: str,
        influxdbToken: str,
        influxdbUrl: str,
    ):
        self.bucket = influxdbBucket
        self.org = influxdbOrg
        self.token = influxdbToken
        self.url = influxdbUrl
