package main

import (
	"api"
	"fmt"
	"logger"
	"time"
)

const tag string = "MAIN"
var Debug string = "true"
var LogDir string
const port string = ":8080"


func main() {
	logger.SetLogFile(LogDir, Debug)
	logger.LogInfo(fmt.Sprintf("Starting server on port %s", port), tag)
	go api.InitServer(port)
	for true {
		logger.LogInfo(fmt.Sprintf("Running on localhost%s", port), tag)
		time.Sleep(time.Hour * 1)
	}
}
