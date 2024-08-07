

class InfluxDBCredentials:
    def __init__(self, influxdbBucket, influxdbOrg, influxdbToken, influxdbUrl):
        self.bucket = influxdbBucket
        self.org = influxdbOrg
        self.token = influxdbToken
        self.url = influxdbUrl