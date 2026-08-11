import threading
import time

from pkg.domain import ImporterAdapter
from pkg.domain import ExporterAdapter
from pkg.domain import Queue
from pkg.converter.otel_converter import Converter
from datetime import datetime
import numpy as np

import logging
log = logging.getLogger(__name__)


# make past data replay in real time
class ReplayTelemetryUseCase:
    def __init__(self, importer: ImporterAdapter, exporter: ExporterAdapter,converter: Converter, queue: Queue, num_workers=4):
        self.importer = importer
        self.exporter = exporter
        self.converter = converter

        self.queue = queue
        self.workers = [Worker(queue, exporter) for _ in range(num_workers)]

        for worker in self.workers:
            worker.start()

    def execute(self) -> bool:
        try:
            for data in self.importer: 

                try:
                    otel_data = self.converter.convert(data, target_index = "CAMC_CT0")
                    
                except Exception as e:
                    log.error(f"Error converting data to OTel format: {e}")
                    return False
                
                try:
                    self.queue.enqueue(otel_data)
                except Exception as e:
                    log.error(f"Error enqueueing OTel data: {e}")
                
        except Exception as e:
            for worker in self.workers:
                worker.stop()
            self.exporter.close()
            log.error(f"Error during export: {e}")
            return False
        
        for worker in self.workers:
            worker.stop()
        self.exporter.close()
        return True
    


class Worker(threading.Thread):
    def __init__(self, queue: Queue, exporter: ExporterAdapter):
        super().__init__()
        
        self.queue = queue
        self.exporter = exporter

        self._stop_event = threading.Event()

    def run(self):
        while not self._stop_event.is_set():
            try:
                otel_data = self.queue.dequeue()
                if otel_data is not None:
                    otel_data.event_timelag_min = (self.now() - otel_data.event_timestamp) / np.timedelta64(1, 'm')
                    self.exporter.export(otel_data)
                    log.info(f"Exported OTel data")
                else:
                    time.sleep(2)  
            except Exception as e:
                log.error(f"Error in export worker: {e}")
    
    def stop(self):
        self._stop_event.set()

    def now(self):
        return np.datetime64("now")