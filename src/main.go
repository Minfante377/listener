package main

import (
	"logger"
	"time"
)

const tag string = "MAIN"
const Debug string = "false"
const LogDir string = "logs"


func main() {
	logger.SetLogFile(LogDir, Debug)
	for true {
		logger.LogInfo("Running...", tag)
		time.Sleep(time.Hour * 1)
	}
}
