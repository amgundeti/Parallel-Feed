package feed

import(
	lock "proj1/lock"
	// "encoding/json"
)


//Feed represents a user's twitter feed
// You will add to this interface the implementations as you complete them.
type Feed interface {
	Add(body string, timestamp float64)
	Remove(timestamp float64) bool
	Contains(timestamp float64) bool
	SendFeed(taskID int ) FeedResponse
}

//feed is the internal representation of a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation. You can assume the feed will not have duplicate Posts
type feed struct {
	Start *post // a pointer to the beginning Post
	locking *lock.LockStruct
}

//Post is the internal representation of a Post on a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation.
type post struct {
	body      string // the text of the Post
	timestamp float64  // Unix timestamp of the Post
	next      *post  // the next Post in the feed
}


type FeedResponse struct{
	Id int `json:"id"`
	Feed []map[string]interface{} `json:"feed"`
}

//NewPost creates and returns a new Post value given its body and timestamp
func newPost(body string, timestamp float64, next *post) *post {
	return &post{body, timestamp, next}
}

//NewFeed creates a empy user feed
func NewFeed() Feed {
	f :=  &feed{Start: nil}
	f.locking = lock.NewLock()
	return f
}

// Add inserts a new Post to the feed. The feed is always ordered by the timestamp where
// the most recent timestamp is at the beginning of the feed followed by the second most
// recent timestamp, etc. You may need to insert a new Post somewhere in the feed because
// the given timestamp may not be the most recent.

func (f *feed) Add(body string, timestamp float64) {
	f.locking.Lock()
	defer f.locking.Unlock()

	if f.Start == nil{
		f.Start = newPost(body, timestamp, nil)
		return
	}

	if f.Start.timestamp <= timestamp{
		newPost := newPost(body, timestamp, f.Start)
		f.Start = newPost
		return
	}

	prev := f.Start
	curr := f.Start.next

	for curr != nil && curr.timestamp > timestamp {
		prev = curr
		curr = curr.next
	}

	newPost := newPost(body, timestamp, curr)
	prev.next = newPost
}

// Remove deletes the Post with the given timestamp. If the timestamp
// is not included in a Post of the feed then the feed remains
// unchanged. Return true if the deletion was a success, otherwise return false
func (f *feed) Remove(timestamp float64) bool {

	f.locking.Lock()
	defer f.locking.Unlock()

	if f.Start == nil{
		return false
	}

	if f.Start.timestamp == timestamp{
		f.Start = f.Start.next
		return true
	}

	prev := f.Start
	curr := f.Start.next

	for curr != nil && curr.timestamp > timestamp{

		prev = curr
		curr = curr.next
	}

	if  curr != nil && curr.timestamp == timestamp{
		prev.next = curr.next
		return true
	}
	return false
}

// Contains determines whether a Post with the given timestamp is
// inside a feed. The function returns true if there is a Post
// with the timestamp, otherwise, false.
func (f *feed) Contains(timestamp float64) bool {

	f.locking.Rlock()
	defer f.locking.RUnlock()

	if f.Start == nil{
		return false
	}

	curr := f.Start

	if curr.timestamp == timestamp{
		return true
	}

	for curr != nil && curr.timestamp >= timestamp {
		if curr.timestamp == timestamp {
			return true
		}
		curr = curr.next
		if curr == nil{
			return false
		}
	}
	return false
}

func (f *feed) SendFeed(taskID int) FeedResponse {

	tempPost := f.Start
	feedResponse := FeedResponse{Id: taskID}

	for tempPost != nil {
		s := make(map[string]interface{})
		s["body"] = tempPost.body
		s["timestamp"] = tempPost.timestamp
		feedResponse.Feed = append(feedResponse.Feed, s)
		tempPost = tempPost.next
	}
	return feedResponse
}

