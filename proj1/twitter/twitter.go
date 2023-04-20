package main

import(
	"os"
	"strconv"
	"proj1/server"
	"encoding/json"
)

func main() {

	configInstance := server.Config{}

	if len(os.Args) == 2{
		configInstance.ConsumersCount, _ = strconv.Atoi(os.Args[1])
	} else{
		configInstance.ConsumersCount = 1
	}


	if configInstance.ConsumersCount > 1 {
		configInstance.Mode = "p"
		configInstance.Decoder = json.NewDecoder(os.Stdin)
		configInstance.Encoder = json.NewEncoder(os.Stdout)
		server.Run(configInstance)
	} else{
		configInstance.Mode = "s"
		configInstance.ConsumersCount = 1
		configInstance.Decoder = json.NewDecoder(os.Stdin)
		configInstance.Encoder = json.NewEncoder(os.Stdout)
		server.Run(configInstance)
	}

}