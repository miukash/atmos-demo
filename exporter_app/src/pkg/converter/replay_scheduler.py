import time

class ReplayScheduler:
    def __init__(self, now_func=None):
        self.gap = None
    def  wait(self, otel_data):
        event_time = otel_data.event_timestamp

        if self.gap is None:
            self.gap = time.time() - event_time
        
        otel_data.event_timestamp = otel_data.event_timestamp + self.gap

