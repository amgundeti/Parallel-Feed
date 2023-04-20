package server

import(
	"encoding/json"
	"sync"
	"proj1/queue"
	"proj1/feed"
	"fmt"

)

type Config struct {
	Encoder *json.Encoder // Represents the buffer to encode Responses
	Decoder *json.Decoder // Represents the buffer to decode Requests
	Mode    string        // Represents whether the server should execute
	// sequentially or in parallel
	// If Mode == "s"  then run the sequential version
	// If Mode == "p"  then run the parallel version
	// These are the only values for Version
	ConsumersCount int // Represents the number of consumers to spawn
}

type Aux struct{
	mu  sync.Mutex
	cond *sync.Cond
	tasks int
	done bool
}

// type Response struct{
// 	Success bool `json:"success"`
// 	Id int	`json:"id"`
// }
type Response struct {
	Success bool  `json:"success"`
	Id      int `json:"id"`
}

//Run starts up the twitter server based on the configuration
//information provided and only returns when the server is fully
// shutdown.
func Run(config Config) {

	aux := &Aux{done: false, tasks: 0}
	aux.cond = sync.NewCond(&aux.mu)
	q := queue.NewLockFreeQueue()
	f := feed.NewFeed()
	var wg sync.WaitGroup

	if config.Mode == "p"{
	for i:= 0; i < config.ConsumersCount; i++{
		wg.Add(1)
		go consumer(aux, q, f, config, &wg)
	}
	producer(aux, q, f, config)
	wg.Wait()
	
	} else{
		sequential(q,f, config)
	}
}

func consumer(aux *Aux, q *queue.LockFreeQueue, f feed.Feed, config Config, wg *sync.WaitGroup){

	var response Response
	var feedResponse feed.FeedResponse
	
	for{
		aux.mu.Lock()
		for aux.tasks == 0 && !aux.done {
			aux.cond.Wait()
		}

		// Check if done command has arrived yet
		if aux.tasks == 0 && aux.done {
			wg.Done()
			aux.mu.Unlock()
			return
		}
		// decreement tasks and go to dequeue
		aux.tasks -= 1
		aux.mu.Unlock()

		task := q.Dequeue()
		response.Id = task.Id

		switch task.Command{
		case "ADD":
			f.Add(task.Body, task.Timestamp)
			response.Success = true
		case "REMOVE":
			response.Success = f.Remove(task.Timestamp)
		case "CONTAINS":
			response.Success = f.Contains(task.Timestamp)
		case "FEED":
			feedResponse = f.SendFeed(task.Id)
		}
	

		if task.Command != "FEED"{
			config.Encoder.Encode(&response)
		} else{
			config.Encoder.Encode(&feedResponse)
		}

	}
}


func producer(aux *Aux, q *queue.LockFreeQueue, f feed.Feed, config Config) {

	for {
		req := &queue.Request{}
		config.Decoder.Decode(req)

		if req.Command == "DONE"{
			aux.mu.Lock()
			aux.done = true
			aux.cond.Broadcast()
			aux.mu.Unlock()
			return
		}

		q.Enqueue(req)
		aux.mu.Lock()
		aux.tasks +=1
		aux.cond.Broadcast()
		aux.mu.Unlock()

	}
}


func sequential(q *queue.LockFreeQueue, f feed.Feed, config Config){

	var request queue.Request
	var response Response
	var feedResponse feed.FeedResponse

	for {

		err := config.Decoder.Decode(&request)
		if err != nil{
			// break
		}

		response.Id = request.Id

		if request.Command == "DONE"{
			fmt.Println("DONE")
			return
		}

		switch request.Command {
		case "ADD":
			f.Add(request.Body, request.Timestamp)
			response.Success = true
		case "REMOVE":
			response.Success = f.Remove(request.Timestamp)
		case "CONTAINS":
			response.Success = f.Contains(request.Timestamp)
		case "FEED":
			feedResponse = f.SendFeed(request.Id)
		}

		if request.Command != "FEED"{
			config.Encoder.Encode(&response)
		} else{
			config.Encoder.Encode(&feedResponse)
		}

	}

}