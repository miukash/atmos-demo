from queue import Queue
from pkg.infra.pkl.npz_importer_ns import NpzImporter
from pkg.infra.matplot.matplot_exporter import MatplotlibExporter
from pkg.queue.queue import RequestQueue
from pkg.usecase.replay_telemetry_data import ReplayTelemetryUseCase
from pkg.converter.otel_converter_ns import NoiseSpectrumConverter


if __name__ == "__main__":

    data_path = "../data/raw_data_jan"
    # log_file_path = "../data/logs/promtail.log"

    usecase = ReplayTelemetryUseCase(
        importer=NpzImporter(
            data_path=data_path
        ),
        exporter=MatplotlibExporter(
            # log_file_path=log_file_path
        ),
        converter=NoiseSpectrumConverter(
            time_window_min=5
        ),
        queue=RequestQueue(size=100)
    )

    usecase.execute()