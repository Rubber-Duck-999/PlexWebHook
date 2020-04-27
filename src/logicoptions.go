package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func messageFailure(issue bool) string {
	fail := ""
	if issue {
		fail = PublishEventNAC(COMPONENT, SERVERERROR, getTime())
	}
	return fail
}

func checkState() {
	for message_id := range SubscribedMessagesMap {
		if SubscribedMessagesMap[message_id].valid == true {
			log.Debug("Message id is: ", message_id)
			log.Debug("Message routing key is: ", SubscribedMessagesMap[message_id].routing_key)
			switch {
			case SubscribedMessagesMap[message_id].routing_key == REQUESTACCESS:
				var message RequestAccess
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				var result string
				log.Debug(message.Pin)
				if message.Pin == pinCode {
					log.Debug("Pins match")
					result = ACCESSPASS
				} else {
					result = ACCESSFAIL
				}
				valid := PublishAccessResponse(message.Id, result)
				if valid != "" {
					log.Warn("Failed to publish")
				} else {
					log.Debug("Published Access Response")
				}
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == DATAINFO:
				log.Warn("Received a data info topic")

			default:
				log.Warn("We were not expecting this message unvalidating: ",
					SubscribedMessagesMap[message_id].routing_key)
				SubscribedMessagesMap[message_id].valid = false
			}
		}
	}

}
