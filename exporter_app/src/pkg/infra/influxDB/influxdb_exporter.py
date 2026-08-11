from influxdb_client import InfluxDBClient, Point
from influxdb_client.client.write_api import ASYNCHRONOUS
from datetime import timezone, timedelta, datetime
import numpy as np

import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)
JST = timezone(timedelta(hours=9))

class InfluxDBExporter:
    def __init__(self, url, token, org, bucket):
        self.client = InfluxDBClient(url=url, token=token, org=org)
        self.write_api = self.client.write_api(write_options=ASYNCHRONOUS)
        self.bucket = bucket
        self.org = org

    def _to_datetime(self, value):
        if isinstance(value, datetime):
            return value

        if isinstance(value, np.datetime64):
            return value.astype('datetime64[s]').astype(datetime)

        if hasattr(value, 'to_pydatetime'):
            return value.to_pydatetime()

        if isinstance(value, str):
            return datetime.fromisoformat(value.replace('Z', '+00:00'))

        raise TypeError(f"Unsupported timestamp type: {type(value)}")

    def export(self, otel_data):
        timestamp = self._to_datetime(otel_data.event_timestamp)
        if timestamp.tzinfo is None:
            timestamp = timestamp.replace(tzinfo=JST)
        timestamp = timestamp.astimezone(timezone.utc)
        
        point = (
            Point("telemetry")
            .field("timelag_min", str(otel_data.event_timelag_min))
            .field("value", float(otel_data.value))
            .field("average", float(otel_data.average))
            .field("std", float(otel_data.std))
            .time(timestamp)
        )

        if otel_data.noise_spectrum is not None:
            point.field("noise_spectrum", str(otel_data.noise_spectrum))

        self.write_api.write(
            bucket=self.bucket,
            org=self.org,
            record=point
        )

        logger.info(f"Exported data to InfluxDB: {otel_data.event_timestamp}, {otel_data.value}")
    
    def close(self):
        self.client.close()