package main

type ConfigTypes struct {
	Settings struct {
		Code    string `yaml:"code"`
	} `yaml:"settings"`
}

type FailureNetwork struct {
	Time         string `json:"time"`
	Failure_type string `json:"type_of_failure"`
}

type RequestImage struct {
	Request_id string `json:"request_id"`
	Time_from string `json:"time_from"`
	Time_to string `json:"time_to"`
}

type RequestData struct {
	Request_id string `json:"request_id"`
	Time_from string `json:"time_from"`
	Time_to string `json:"time_to"`
	Type string `json:"type"`
}

type RequestDatabase struct {
	Request_id string `json:"request_id"`
	Time_from string `json:"time_from"`
	Time_to string `json:"time_to"`
	Type string `json:"type"`
}

type DataInfo struct {
	Id int `json:"id"`
    Message string `json:"message"`
	Time string `json:"Time"`
}

type StatusUpdate struct {
    Cpu_usage int `json:"cpu_usage"`
    Memory_used int `json:"memory_used"`
    Images int `json:"images"`
}

type RequestAccess struct {
	Id int `json:"id"`
	Pin int `json:"pin"`
}

type EventNAC struct {
	Component    string
	Message      string
	Time         string
}

type MapMessage struct {
	message     string
	routing_key string
	time        string
	valid       bool
}

//Topics
const REQUESTIMAGE string = "Request.Image"
const REQUESTDATA string = "Request.Data"
const AUTHENTICATIONREQUEST string = "Authentication.Request"
const DATAINFO string = "Data.Info"
const STATUSUPDATE string = "Status.Update"
const REQUESTACCESS string = "Request.Access"
//
const FAILURENETWORK string = "Failure.Network"
const EVENTNAC string = "Event.NAC"
const REQUESTDATABASE string = "Request.Database"
const DATARESPONSE string = "Data.Response"
const DEVICEFOUND string = "Device.Found"
const AUTHENTICATIONRESPONSE string = "Authentication.Response"
const ACCESSRESPONSE string = "Access.Response"

const EXCHANGENAME string = "topics"
const EXCHANGETYPE string = "topic"
const TIMEFORMAT string = "20060102150405"
const COMPONENT string = "NAC"
const FAILUREPUBLISH string = "Failed to publish"
const SERVERERROR string = "Server is failing to send"

var SubscribedMessagesMap map[uint32]*MapMessage
var key_id uint32 = 0
