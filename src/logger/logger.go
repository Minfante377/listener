package logger

import(
	"fmt"
	"log"
	"os"
	"time"
)

const(
	layout string = "01-02-2006_03:04"
)

var debug_mode string = ""
var log_chan chan string = make(chan string, 100)


func write() {
	log.Print(<-log_chan)
}


func SetLogFile(filepath string, debug string) int {
	os.Mkdir(filepath, 0777)
	dt := time.Now().Format("01-02-2006")
	log_name := fmt.Sprintf("%s/%s.log", filepath, dt)
	file, err := os.OpenFile(log_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY,
							 0666)
	if err != nil {
		return -1
	}
	log.SetOutput(file)
	debug_mode = debug
	return 0
}


func LogInfo(msg string, tag string) {
	var dt string = time.Now().Format(layout)
	var log_msg string = fmt.Sprintf("[%s - INFO - %s]: %s\n", tag, dt, msg)
	log_chan <- log_msg
	go write()
	if debug_mode == "true" {
		fmt.Printf("[%s- INFO - %s]: %s\n", tag, dt, msg)
	}
}


func LogError(msg string, tag string) {
	var dt string = time.Now().Format(layout)
	var log_msg string = fmt.Sprintf("[%s - ERROR - %s]: %s\n", tag, dt, msg)
	log_chan <- log_msg
	go write()
	if debug_mode == "true" {
		fmt.Printf("[%s - ERROR - %s]: %s\n", tag, dt, msg)
	}
}


func LogTestStep(step string) {
	log.Printf(fmt.Sprintf("[TEST STEP]: %s---------------\n\n", step))
	if debug_mode == "true" {
		fmt.Printf(fmt.Sprintf("[TEST STEP]: %s---------------\n\n", step))
	}
}
