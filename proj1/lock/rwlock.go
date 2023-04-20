package lock

import(
	"sync"
)

type LockStruct struct {
	mu sync.Mutex
	cond *sync.Cond
	readersMax int64
	readersActive int64
	writersWaiting int64
	writerActive bool
}

func NewLock() *LockStruct{
	r := &LockStruct{readersMax: 32, readersActive: 0, writersWaiting: 0, writerActive: false}
	r.cond = sync.NewCond(&r.mu)
	return r
}

func (r * LockStruct) Lock(){

	r.mu.Lock()
	defer r.mu.Unlock()
	
	for r.writerActive  || r.readersActive > 0 {
		r.cond.Wait()
	}

	r.writerActive = true
}

func (r * LockStruct) Unlock(){

	r.mu.Lock()
	defer r.mu.Unlock()

	r.writerActive = false
	r.cond.Signal()
}

func (r * LockStruct) Rlock(){

	r.mu.Lock()
	defer r.mu.Unlock()

	for r.readersActive > r.readersMax || r.writerActive {
		r.cond.Wait()
	}

	r.readersActive += 1
}

func (r * LockStruct) RUnlock(){
	r.mu.Lock()
	defer r.mu.Unlock()
	r.readersActive -= 1
	if r.readersActive == 0{
		r.cond.Signal()
	}
	
}
////////////////////////////////////////////////////////////////////////////////////////
