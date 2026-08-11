from queue import Queue
from pkg.infra.influxDB.influxdb_exporter import InfluxDBExporter
from pkg.infra.pkl.pkl_importer import PklImporter
from pkg.infra.promtail.promtail_exporter import PromtailExporter
from pkg.queue.queue import RequestQueue
from pkg.usecase.replay_telemetry_data import ReplayTelemetryUseCase
from pkg.converter.otel_converter import Converter


if __name__ == "__main__":

    data_path = "../data/raw_data_jan"
    # log_file_path = "../data/logs/promtail.log"

    usecase = ReplayTelemetryUseCase(
        importer=PklImporter(
            pkl_data_path=data_path
        ),
        exporter=InfluxDBExporter(
            url="http://localhost:8086",
            token="yiWQres-JdWOZaNxSyd-UhVdA3hDh0JWyCOuiFYdg1lnzwyGTkXet5s_m-atPFMU5LyhzwCX5hugvB-zZNgYhw==",
            org="personal",
            bucket="personal"
        ),
        converter=Converter(
            time_window_min=5
        ),
        queue=RequestQueue(size=100)
    )

    usecase.execute()