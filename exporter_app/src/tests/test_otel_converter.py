from pkg.converter.otel_converter import Converter
import numpy as np

class TestConverter:
    def test_convert(self):
        
        event = {
            "TimeStamp": np.datetime64("2024-05-01T15:00:00Z"),
            "value": 40
        }

        past_timestamps = [
            np.datetime64('2024-05-01T12:00:00Z'),
            np.datetime64('2024-05-01T13:00:00Z'),
            np.datetime64('2024-05-01T14:00:00Z')
        ] 
        past_values = [10, 20, 30]


        converter = Converter(time_window_min=150)
        converter.timestamps = past_timestamps[:]
        converter.values = past_values[:]

        otel_data = converter.convert(event, target_index="value")
        expected_values = past_values[1:]+[event["value"]]

        assert otel_data.event_timestamp == event["TimeStamp"]
        assert otel_data.value == event["value"]
        assert otel_data.average == np.mean(expected_values)
        assert otel_data.std == np.std(expected_values)

    # def test_convert_with_noise_spectrum(self):
    #     event = {
    #         "TimeStamp": np.datetime64("2024-05-01T15:00:00Z"),
    #         "value": 40,
    #         "NoiseSpectrum": [1.0, 2.0, 3.0]
    #     }

    #     converter = Converter(time_window_min=150)
    #     otel_data = converter.convert(event, target_index="value")

    #     assert otel_data.value == event["value"]
    #     assert otel_data.noise_spectrum == [1.0, 2.0, 3.0]