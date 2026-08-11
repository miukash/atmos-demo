import threading
from unittest.mock import MagicMock, patch
from pkg.usecase.replay_telemetry_data import ReplayTelemetryUseCase, Worker
from pkg import domain
import numpy as np


class TestReplayTelemetryUseCase:

    @patch('pkg.usecase.replay_telemetry_data.Worker')
    def test_execute(self, mock_worker):

        mock_importer = MagicMock()
        mock_importer.__iter__.return_value = iter([
            {'timestamp': 1, 'value': 10},
            {'timestamp': 2, 'value': 20},
            {'timestamp': 3, 'value': 30}
        ])
        mock_exporter = MagicMock()
        
        mock_queue = MagicMock()

        mock_converter = MagicMock()
        mock_converter.convert_to_otel.return_value = domain.OtelData(0, 0, 0, 0, 0)

        
        exporter = ReplayTelemetryUseCase(mock_importer, mock_exporter, mock_converter, mock_queue)
        mock_worker_instance = mock_worker.return_value
        mock_worker_instance.start = MagicMock()

        assert exporter.execute() == True

    @patch('pkg.usecase.replay_telemetry_data.Worker')
    def test_execute_with_importer_error(self, mock_worker):
        mock_importer = MagicMock()
        mock_importer.__iter__.side_effect = Exception("Importer error")
        mock_exporter = MagicMock()
        mock_queue = MagicMock()
        mock_converter = MagicMock()

        exporter = ReplayTelemetryUseCase(mock_importer, mock_exporter, mock_converter, mock_queue)
        mock_worker_instance = mock_worker.return_value
        mock_worker_instance.start = MagicMock()

        assert exporter.execute() == False

    @patch('pkg.usecase.replay_telemetry_data.Worker.now')
    def test_export_worker(self, mock_np_datetime64_now):

        enqueued_data = domain.OtelData(
            event_timestamp=np.datetime64("2024-05-01T15:00:00Z"),
            event_timelag_min=0,
            value=10,
            average=10,
            std=0
        )
        mock_np_datetime64_now.return_value = np.datetime64("2024-05-01T15:05:00Z")
        
        
        mock_importer = MagicMock()
        mock_exporter = MagicMock()
        mock_exporter.export = MagicMock()

        mock_queue = MagicMock()
        mock_queue.dequeue.return_value = enqueued_data
        
        mock_converter = MagicMock()

        exporter = ReplayTelemetryUseCase(mock_importer, mock_exporter, mock_converter, mock_queue)


        t = threading.Thread(target=exporter.workers[0].start)
        t.start()
        stop_event = threading.Event()
        stop_event.wait(timeout=1)
        for worker in exporter.workers:
            worker._stop_event.set()
        t.join(timeout=1)   

        assert mock_queue.dequeue.called
        assert mock_exporter.export.called
        assert enqueued_data.event_timelag_min == 5.0





    