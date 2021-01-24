package main

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
)

func convertStatus(status string) int {
	switch {
	case strings.Contains(status, ALLOWED_STRING):
		return ALLOWED
	case strings.Contains(status, BLOCKED_STRING):
		return BLOCKED
	case strings.Contains(status, UNKNOWN_STRING):
		return UNKNOWN
	default:
		return DISCOVERED
	}
}

func deviceResponse(request_id uint32) {
	if DevicesList[request_id].Allowed == BLOCKED {
		PublishDeviceFound(DevicesList[request_id].Device_name,
			DevicesList[request_id].Ip_address,
			DevicesList[request_id].Allowed)
	} else if DevicesList[request_id].Allowed == DISCOVERED {
		log.Error("We did not get a correct status")
	} else if DevicesList[request_id].Allowed == ALLOWED {
		log.Trace("Device is allowed")
	} else if DevicesList[request_id].Allowed == UNKNOWN {
		PublishDeviceFound(DevicesList[request_id].Device_name,
			DevicesList[request_id].Ip_address,
			DevicesList[request_id].Allowed)
	} else {
		log.Error("We shouldn't hit this error")
	}
}

func checkState() {
	for message_id := range SubscribedMessagesMap {
		if SubscribedMessagesMap[message_id].valid {
			switch {
			case SubscribedMessagesMap[message_id].routing_key == DEVICERESPONSE:
				log.Warn("Received a device response topic")
				var message DeviceResponse
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				DevicesList[message.Request_id].Allowed = convertStatus(message.Status)
				DevicesList[message.Request_id].Device_name = message.Name
				deviceResponse(message.Request_id)
				SubscribedMessagesMap[message_id].valid = false

			default:
				log.Warn("We were not expecting this message unvalidating: ",
					SubscribedMessagesMap[message_id].routing_key)
				SubscribedMessagesMap[message_id].valid = false
			}
		}
	}

}
