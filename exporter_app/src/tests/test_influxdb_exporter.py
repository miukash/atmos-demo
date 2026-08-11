from pkg.domain import OtelData
from pkg.infra.influxDB.influxdb_exporter import InfluxDBExporter
from unittest.mock import patch

class TestInfluxDBExporter:
    @patch("pkg.infra.influxDB.influxdb_exporter.InfluxDBClient")
    def test_export(self, mock_influx_client):
        mock_write_api = mock_influx_client.return_value.write_api.return_value

        exporter = InfluxDBExporter(
            url="http://localhost:8086",
            token="my-token",
            org="my-org",
            bucket="my-bucket"
        )

        otel_data = OtelData(
            event_timestamp="2024-01-01T00:00:00Z",
            event_timelag_min=5,
            value=42.0,
            average=40.0,
            std=2.0
        )

        exporter.export(otel_data)

        mock_write_api.write.assert_called_once()

        _, kwargs = mock_write_api.write.call_args

        assert kwargs["bucket"] == "my-bucket"
        assert kwargs["org"] == "my-org"

        point = kwargs["record"]

        assert point._name == "telemetry"
        assert point._fields["timelag_min"] == "5"
        assert point._fields["value"] == 42.0
        assert point._fields["average"] == 40.0
        assert point._fields["std"] == 2.0

    @patch("pkg.infra.influxDB.influxdb_exporter.InfluxDBClient")
    def test_export_with_pandas_timestamp(self, mock_influx_client):
        mock_write_api = mock_influx_client.return_value.write_api.return_value

        exporter = InfluxDBExporter(
            url="http://localhost:8086",
            token="my-token",
            org="my-org",
            bucket="my-bucket"
        )

        import pandas as pd

        otel_data = OtelData(
            event_timestamp=pd.Timestamp("2024-01-01T00:00:00Z"),
            event_timelag_min=5,
            value=42.0,
            average=40.0,
            std=2.0
        )

        exporter.export(otel_data)

        mock_write_api.write.assert_called_once()

    @patch("pkg.infra.influxDB.influxdb_exporter.InfluxDBClient")
    def test_export_with_noise_spectrum(self, mock_influx_client):
        mock_write_api = mock_influx_client.return_value.write_api.return_value

        exporter = InfluxDBExporter(
            url="http://localhost:8086",
            token="my-token",
            org="my-org",
            bucket="my-bucket"
        )

        otel_data = OtelData(
            event_timestamp="2024-01-01T00:00:00Z",
            event_timelag_min=5,
            value=42.0,
            average=40.0,
            std=2.0,
            noise_spectrum=[1.0, 2.0, 3.0]
        )

        exporter.export(otel_data)

        mock_write_api.write.assert_called_once()

        _, kwargs = mock_write_api.write.call_args
        point = kwargs["record"]

        assert point._fields["noise_spectrum"] == "[1.0, 2.0, 3.0]"