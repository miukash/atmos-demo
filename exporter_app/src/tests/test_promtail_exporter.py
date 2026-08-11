import tempfile
import os
from pkg.domain import OtelData
import numpy as np

from pkg.infra.promtail.promtail_exporter import PromtailExporter

class TestPromtailExporter:
    def test_export(self):
        
        # Create a temporary directory to store the log file
        with tempfile.TemporaryDirectory() as temp_dir:
            log_file_path = os.path.join(temp_dir, "test.log")
            expected_data = OtelData(
                event_timestamp=np.datetime64('2024-01-01T00:00:00'),
                event_timelag_min=5.0,
                value=42.0,
                average=40.0,
                std=2.0
            )
            
            with open(log_file_path, 'w') as log_file:
                exporter = PromtailExporter(log_file_path)
                exporter.export(expected_data)
                exporter.close()

                assert os.path.exists(log_file_path), "Log file should exist after export"
                with open(log_file_path, 'r') as log_file:
                    log_content = log_file.read()
                    assert str(expected_data.event_timestamp) in log_content, "Event timestamp should be in log content"
                    assert str(expected_data.event_timelag_min) in log_content, "Event timelag should be in log content"
                    assert str(expected_data.value) in log_content, "Value should be in log content"
                    assert str(expected_data.average) in log_content, "Average should be in log content"
                    assert str(expected_data.std) in log_content, "Std should be in log content"
            
            
