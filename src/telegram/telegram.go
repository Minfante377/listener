package telegram

import (
	"api"
	"bufio"
	"fmt"
	"logger"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const tag string = "TELEGRAM"
const sendMsgApi string = "https://api.telegram.org/bot%s/"+
						  "sendMessage?chat_id=%s&text=%s"

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
	var request string = fmt.Sprintf(sendMsgApi, config.botToken, user, msg)
	_, err := http.Get(request)
	if err != nil {
		logger.LogError("Error sending msg", tag)
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
				var msg string = fmt.Sprintf("New critical event %d at %s",
											 event.EventType, event.Date)
				for _, user := range criticalEvent.users {
					sendMsg(msg, user)
				}
			}
		}
	}
}
