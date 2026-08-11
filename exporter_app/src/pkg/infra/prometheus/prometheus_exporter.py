from pkg.domain import ExporterAdapter
from prometheus_client import Gauge
from prometheus_client import start_http_server


class PrometheusExporter(ExporterAdapter):
    def __init__(self):
        self.timelag = Gauge("telemetry_timelag_min", "Telemetry timelag in minutes")
        self.value = Gauge("telemetry_value", "Telemetry value")
        self.avg = Gauge("telemetry_avg", "Telemetry average")
        self.std = Gauge("telemetry_std", "Telemetry standard deviation")

        start_http_server(8000)

    def export_data(self, data):
        self.timelag.set(data.event_timelag_min)
        self.value.set(data.value)
        self.avg.set(data.average)
        self.std.set(data.std)
