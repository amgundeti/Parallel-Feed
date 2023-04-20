package queue

import (
	"sync/atomic"
	"unsafe"
)

type Request struct {
	Command string
	Id int
	Body string
	Timestamp float64
	Next *Request
}


type LockFreeQueue struct {
	head  *Request
	tail  *Request
}

//Attribution: https://www.sobyte.net/post/2021-07/implementing-lock-free-queues-with-go/ && Herlihy & Shavit

func NewLockFreeQueue() *LockFreeQueue {
	sentinel := &Request{}
	return &LockFreeQueue{head: sentinel, tail: sentinel}
}	


func (queue *LockFreeQueue) Enqueue(task *Request) {
	
	last := queue.tail
	for !cas(&last.Next, nil, task){
		// find tail again
		last = queue.tail
	}
	cas(&queue.tail, last, task)
}


func (queue *LockFreeQueue) Dequeue() *Request {

	head := queue.head

	// check if queue empty
	if head.Next == nil {
		return nil
	}

	for {
		// if successfull in moving head along to next node return
		if cas(&queue.head, head, head.Next){
			return head.Next
		}

		// otherwise find head again
		head = queue.head
		
		// make sure queue isn't empty
		if(head.Next == nil) {
			return nil
		}
	}
}

func cas(old **Request, expected *Request, new *Request) bool {
	return atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(old)), (unsafe.Pointer)(expected), (unsafe.Pointer)(new))
}
