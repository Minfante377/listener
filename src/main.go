package main

import (
	"api"
	"fmt"
	"logger"
	"time"
	"telegram"
)

const tag string = "MAIN"
var Debug string = "true"
var LogDir string
const port string = ":8080"
const config_path string = "./bot.config"


func main() {
	logger.SetLogFile(LogDir, Debug)
	logger.LogInfo(fmt.Sprintf("Starting server on port %s", port), tag)
	var critical_events chan api.Event = make(chan api.Event, 10)
	go api.InitServer(port, critical_events)
	go telegram.BotHandler(critical_events, config_path)
	for true {
		logger.LogInfo(fmt.Sprintf("Running on localhost%s", port), tag)
		time.Sleep(time.Hour * 1)
	}
}
