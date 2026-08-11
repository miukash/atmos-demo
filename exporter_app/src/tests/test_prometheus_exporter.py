from unittest.mock import Mock, patch
from pkg.infra.prometheus.prometheus_exporter import PrometheusExporter
from pkg.domain import OtelData

class TestPrometheusExporter:

    @patch('prometheus_client.start_http_server')
    def test_export_data(self, mock_start_http_server):
        mock_start_http_server.return_value = None
        
        exporter = PrometheusExporter()

        exporter.timelag.set = Mock()
        exporter.value.set = Mock()
        exporter.avg.set = Mock()
        exporter.std.set = Mock()
        
        data = OtelData(
            event_timestamp='2024-01-01T00:00:00Z',
            event_timelag_min=5, 
            value=100, average=95, std=10)

        exporter.export_data(data)

        exporter.timelag.set.assert_called_once_with(5)
        exporter.value.set.assert_called_once_with(100)
        exporter.avg.set.assert_called_once_with(95)
        exporter.std.set.assert_called_once_with(10)
