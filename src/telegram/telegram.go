package telegram

import (
	"api"
	"bytes"
	"bufio"
	"encoding/json"
	"fmt"
	"logger"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const tag string = "TELEGRAM"
const sendMsgApi string = "https://api.telegram.org/bot%s/sendMessage?"

var config configuration = configuration{}

type criticalEvent struct {
	eventType int32
	users []string
}

type configuration struct {
	botToken string
	criticalEvents []criticalEvent
}


func readConfig(path string) (configuration, int) {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error reading config file: %s",
									err.Error()), tag)
		return config, 1
	}
	scanner := bufio.NewScanner(f)
	var bot_token string = ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "//") || len(line) == 0 {
			continue
		}else if strings.Contains(line, "token") {
			bot_token = strings.Split(line, "=")[1]
		}else {
			split_line := strings.Split(line, "=")
			event_str := split_line[0]
			users := split_line[1]
			var critical_event criticalEvent
			event, _ := strconv.Atoi(event_str)
			critical_event.eventType = int32(event)
			for _, user := range strings.Split(users, ",") {
				critical_event.users = append(critical_event.users, user)
			}
			config.criticalEvents = append(config.criticalEvents,
										   critical_event)
		}
	}
	if len(bot_token) == 0 {
		logger.LogError("Missing TELEGRAM TOKEN on config file", tag)
		return config, 1
	}
	config.botToken = bot_token
	return config, 0
}


func sendMsg(msg string, user string) {
	logger.LogInfo(fmt.Sprintf("Sending message %s to user %s", msg, user),
				   tag)
	var url string = fmt.Sprintf(sendMsgApi, config.botToken)
	post_body, _ := json.Marshal(map[string]string{
		"chat_id":  user,
		"text": msg,
	})
	response_body := bytes.NewBuffer(post_body)
	_, err := http.Post(url, "Application/json", response_body)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error sending msg: %s", err.Error()), tag)
	}
}


func BotHandler(critical_events chan api.Event, config_path string) {
	config, err := readConfig(config_path)
	if err != 0 {
		panic("Error reading configuration!")
	}
	for true {
		event := <-critical_events
		for _, criticalEvent := range config.criticalEvents {
			if event.EventType == criticalEvent.eventType {
				var msg string = fmt.Sprintf("New critical event %d at %s.\n"+
											 "User: %s\nPwd: %s\nCmd: %s"+
											 "\nPid: %s\nNotes: %s",
											 event.EventType, event.Date,
									         event.User, event.Pwd, event.Cmd,
										     event.Pid, event.Notes)
				for _, user := range criticalEvent.users {
					sendMsg(msg, user)
				}
			}
		}
	}
}
