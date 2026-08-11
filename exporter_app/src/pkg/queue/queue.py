from pkg.domain import Queue
import queue

class RequestQueue(Queue):
    def __init__(self, size=None):
        self.timeout = 10
        maxsize = 0 if size is None else size
        self.queue = queue.Queue(maxsize=maxsize)
    
    def enqueue(self, otel_data):
        self.queue.put(otel_data, timeout=self.timeout)

    def dequeue(self):
        if self.queue:
            try:
                item = self.queue.get(timeout=self.timeout)
                return item
            except queue.Empty:
                return None
        return None
    
    # TODO: if queue is full, wait until queue is empty or timeout
