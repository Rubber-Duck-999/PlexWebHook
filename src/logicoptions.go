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

func convertStatus(status string) int {
	switch {
	case status == ALLOWED_STRING:
		return ALLOWED
	case status == BLOCKED_STRING:
		return BLOCKED
	case status == UNKNOWN_STRING:
		return UNKNOWN
	default:
		return DISCOVERED
	}
}

func deviceResponse(request_id uint32) {
	if DevicesList[request_id].Allowed == BLOCKED  || 
		DevicesList[request_id].Allowed == UNKNOWN {
		PublishDeviceFound(DevicesList[request_id].Device_name,
			DevicesList[request_id].Ip_address,
			DevicesList[request_id].Allowed)
	} else if DevicesList[request_id].Allowed == DISCOVERED {
		log.Error("DBM did not send us a correct status")
	} else if DevicesList[request_id].Allowed == ALLOWED {
		log.Trace("Device is allowed")
	} else {
		log.Error("We shouldn't hit this error")
		log.Error("Allowed status: ", DevicesList[request_id].Allowed)
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
				DevicesList[Request_id].Allowed = convertStatus(message.Status)
				DevicesList[Request_id].Device_name = message.Name
				deviceResponse(Request_id)
				SubscribedMessagesMap[message_id].valid = false
			default:
				log.Warn("We were not expecting this message unvalidating: ",
					SubscribedMessagesMap[message_id].routing_key)
				SubscribedMessagesMap[message_id].valid = false
			}
		}
	}

}
