from pkg.queue.queue import RequestQueue
import queue

class TestRequestQueue:
    def test_enqueue_and_dequeue(self):

        rq = RequestQueue()
        rq.enqueue("request1")
        rq.enqueue("request2")

        assert rq.dequeue() == "request1"
        assert rq.dequeue() == "request2"
        assert rq.dequeue() is None

 
        