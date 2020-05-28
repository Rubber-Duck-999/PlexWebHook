package main

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

func messageFailure(issue bool) string {
	fail := ""
	if issue {
		fail = PublishEventNAC(COMPONENT, SERVERERROR, getTime())
	}
	return fail
}

func checkDay(daily int) int {
	_, _, current_day := time.Now().Date()
	if current_day == day {
		daily++
		return daily
	} else {
		_, _, day = time.Now().Date()
		return 0
	}
}

func convertStatus(status string) int {
	switch {
	case status == ALLOWED_STRING:
		_statusNAC.DailyAllowedDevices = checkDay(_statusNAC.DailyAllowedDevices)
		return ALLOWED
	case status == BLOCKED_STRING:
		_statusNAC.DailyBlockedDevices = checkDay(_statusNAC.DailyBlockedDevices)
		return BLOCKED
	case status == UNKNOWN_STRING:
		_statusNAC.DailyUnknownDevices = checkDay(_statusNAC.DailyUnknownDevices)
		return UNKNOWN
	default:
		return DISCOVERED
	}
}

func checkState() {
	for message_id := range SubscribedMessagesMap {
		if SubscribedMessagesMap[message_id].valid == true {
			log.Debug("Message id is: ", message_id)
			log.Debug("Message routing key is: ", SubscribedMessagesMap[message_id].routing_key)
			switch {
			case SubscribedMessagesMap[message_id].routing_key == DATAINFO:
				log.Warn("Received a data info topic")

			case SubscribedMessagesMap[message_id].routing_key == DEVICERESPONSE:
				log.Warn("Received a device response topic")
				var message DeviceResponse
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				Request_id := message.Request_id
				log.Warn("Device name is: ", message.Name)
				if DevicesList[Request_id].Alive == true {
					DevicesList[Request_id].Allowed = convertStatus(message.Status)
					DevicesList[Request_id].Device_name = message.Name
					if DevicesList[Request_id].Allowed == BLOCKED  || 
						DevicesList[Request_id].Allowed == UNKNOWN {
						PublishDeviceFound(DevicesList[Request_id].Device_name,
							DevicesList[Request_id].Ip_address,
							DevicesList[Request_id].Allowed)
					} else if DevicesList[Request_id].Allowed == DISCOVERED {
						log.Error("DBM did not send us a correct status")
					} else if DevicesList[Request_id].Allowed == ALLOWED {
						log.Trace("Device is allowed")
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
