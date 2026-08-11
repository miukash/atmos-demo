

# infra adapter

from abc import abstractmethod
import numpy as np


class ImporterAdapter:
    @abstractmethod
    def import_data(self):
        # Implement the logic to import data from the source
        pass

class ExporterAdapter:
    @abstractmethod
    def export_data(self, otel_data):
        # Implement the logic to export data to the destination
        pass

# domain entity

class RawData:
    def __init__(self, event_timestamp, value):
        self.event_timestamp : np.datetime64 = event_timestamp
        self.value : float = value

class OtelData:
    def __init__(self, event_timestamp, event_timelag_min, value, average, std, noise_spectrum=None):
        self.event_timestamp : np.datetime64 = event_timestamp
        self.event_timelag_min : float = event_timelag_min
        self.value : float = value
        self.average : float = average
        self.std : float = std
        self.noise_spectrum = noise_spectrum

        

class Queue:
    @abstractmethod
    def enqueue(self, otel_data):
        # Implement the logic to enqueue data for processing
        pass

from dataclasses import dataclass
from datetime import datetime
import numpy as np

@dataclass
class NoiseSpectrumData:
    timestamp: datetime
    pixel: int
    frequency: np.ndarray
    power: np.ndarray
