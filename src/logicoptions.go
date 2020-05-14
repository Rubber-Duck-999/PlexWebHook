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
				if message.Pin == pinCode {
					log.Debug("Pins match")
					result = ACCESSPASS
				} else {
					log.Debug("Pins do not match")
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

			case SubscribedMessagesMap[message_id].routing_key == DEVICERESPONSE:
				log.Warn("Received a device response topic")
				var message DeviceResponse
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				Request_id := message.Request_id
				log.Warn("Device Request for ID: ", Request_id)
				log.Debug("Allowed status: ", DevicesList[Request_id].Allowed, " changing to ",
					message.Status)
				if DevicesList[Request_id].Alive == true {
					if message.Status == ALLOWED_STRING {
						DevicesList[Request_id].Allowed = ALLOWED
					} else if message.Status == BLOCKED_STRING {
						DevicesList[Request_id].Allowed = BLOCKED
					} else {
						DevicesList[Request_id].Allowed = UNKNOWN
					}
					DevicesList[Request_id].Device_name = message.Name
					if DevicesList[Request_id].Allowed == BLOCKED  || 
						DevicesList[Request_id].Allowed == UNKNOWN {
						PublishDeviceFound(DevicesList[Request_id].Device_name,
							DevicesList[Request_id].Ip_address,
							DevicesList[Request_id].Allowed)
					} else if DevicesList[Request_id].Allowed == DISCOVERED {
						log.Error("DBM did not send us a correct status")
					} else if DevicesList[Request_id].Allowed == ALLOWED {
						log.Debug("Device is allowed")
					} else {
						log.Error("We shouldn't hit this error")
						log.Error("Allowed status: ", DevicesList[Request_id].Allowed)
					}
				} else {
					log.Error("We received a response for a non existent device")
				}
				SubscribedMessagesMap[message_id].valid = false
			default:
				log.Warn("We were not expecting this message unvalidating: ",
					SubscribedMessagesMap[message_id].routing_key)
				SubscribedMessagesMap[message_id].valid = false
			}
		}
	}

}
