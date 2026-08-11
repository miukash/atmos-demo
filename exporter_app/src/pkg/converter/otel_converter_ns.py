from pkg.domain import NoiseSpectrumData
import numpy as np


class NoiseSpectrumConverter:
    def __init__(self):

        fs = 12.5 * 10**3
        N = 4096
        self.freq = np.arange(1, N + 1) * (fs / 2 / N)

    def convert(self, event):

        for pixel in range(36):

            spectrum = np.sqrt(event["Spectrum"].T[pixel])

            yield NoiseSpectrumData(
                timestamp=event["TimeStamp"],
                pixel=pixel,
                frequency=self.freq,
                power=spectrum
            )