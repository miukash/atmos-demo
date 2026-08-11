import json
import os


class PromtailExporter:
    def __init__(self, log_file_path):
        self.log_file_path = log_file_path
        self._log_file = None
        
    def _ensure_open(self):
        if self._log_file is None or self._log_file.closed:
            if not os.path.exists(self.log_file_path):
                os.makedirs(os.path.dirname(self.log_file_path), exist_ok=True)
            self._log_file = open(self.log_file_path, 'a')
        
    def export(self, otel_data):
        if hasattr(otel_data.event_timestamp, "isoformat"):
            timestamp_str = otel_data.event_timestamp.isoformat()
        else:
            timestamp_str = str(otel_data.event_timestamp)

        json_log_entry = {
            "event_timestamp": timestamp_str,
            "event_timelag_min": otel_data.event_timelag_min,
            "value": float(otel_data.value) if otel_data.value is not None else None,
            "average": float(otel_data.average) if otel_data.average is not None else None,
            "std": float(otel_data.std) if otel_data.std is not None else None
        }
        self._ensure_open()
        self._log_file.write(json.dumps(json_log_entry) + '\n')

    def close(self):
        if self._log_file and not self._log_file.closed:
            self._log_file.close()

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()