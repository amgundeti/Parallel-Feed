for {
	request := &queue.Request
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Print("Error in decoding request : ", request)
		continue
	}

	if request.Command == "DONE" {
		break
	}

	response := &queue.Response{}

	response.Id = request.Id

	switch request.Command{
	case "ADD":
		f.Add(req.Body, req.Timestamp)
		response.Success = true
	case "REMOVE":
		response.Success = f.Remove(req.Timestamp)
	case "CONTAINS":
		response.Success = f.Contains(req.Timestamp)
	case "FEED":
		feedResponse = f.SendFeed(req.Id)
	}

	err = encoder.Encode(&response)
	if err != nil {
		fmt.Print("Error in encoding response : ", response)
		continue
	}
}