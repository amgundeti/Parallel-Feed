package semaphore

import(
	"sync"
)

type Semaphore struct {
	mu sync.Mutex
	tasks int
	capacity int
	Cond *sync.Cond
	Done bool
	
}

func NewSemaphore() *Semaphore {
	s:= &Semaphore{
		tasks:0,
		capacity: 50,
		Done: false}
	s.Cond = sync.NewCond(&s.mu)
	return s
}
func (s *Semaphore) TaskUp() {
	s.mu.Lock()
	s.tasks += 1
	s.Cond.Broadcast()
	s.mu.Unlock()

}

func (s *Semaphore) CapacityDown() {
	s.mu.Lock()
	for s.capacity == 0 && !s.Done {
		s.Cond.Wait()
	}

	if s.Done{
		s.mu.Unlock()
		return
	}

	s.capacity -= 1
	s.mu.Unlock()

}


func (s *Semaphore) TaskDown() {
	s.mu.Lock()
	
	for s.tasks == 0 && !s.Done{
		s.Cond.Wait()
	}

	if s.Done {
		s.mu.Unlock()
		return
	}
	s.tasks -= 1
	s.mu.Unlock()

}

func (s *Semaphore) CapacityUp() {

	s.mu.Lock()
	s.capacity +=1
	s.Cond.Broadcast()
	s.mu.Unlock()

}
