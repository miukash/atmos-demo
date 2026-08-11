# domain service

from pkg.domain import OtelData
import numpy as np


class Converter: 
    def __init__(self, time_window_min):
        self.time_window_min = time_window_min
        self.timestamps = []
        self.values = []

    def convert(self, event, target_index) -> OtelData:

        self.timestamps.append(event["TimeStamp"])
        self.values.append(event[target_index])

        # Remove old data outside the time window
        while self.timestamps:
            delta = self.timestamps[-1] - self.timestamps[0]
            delta_min = delta / np.timedelta64(1, "m")
            if delta_min > self.time_window_min:
                self.timestamps.pop(0)
                self.values.pop(0)
            else:
                break
        
        if len(self.values) == 0:
            average = None
            std = None
        elif len(self.values) == 1:
            average = float(self.values[0])
            std = 0.0
        else:
            average = float(np.mean(self.values))
            std = float(np.std(self.values))

        noise_spectrum = event.get("NoiseSpectrum")
        if noise_spectrum is not None and not isinstance(noise_spectrum, list):
            noise_spectrum = list(noise_spectrum)

        otel_data = OtelData(
            event_timestamp=event["TimeStamp"],
            event_timelag_min=None,
            value=event[target_index],
            average=average,
            std=std,
        )

        return otel_data
