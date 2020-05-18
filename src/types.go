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
	Request_id int `json:"request_id"`
	Time_from string `json:"time_from"`
	Time_to string `json:"time_to"`
	Type string `json:"type"`
}

type DataInfo struct {
	Id int `json:"id"`
    Message string `json:"message"`
	Time string `json:"Time"`
}

type UnauthorisedConnection struct {
	Mac string `json:"mac"`
	Time string `json:"time"`
	Alive bool `json:"alive"`
}

type EventNAC struct {
	Component    string `json:"Component"`
	Message      string `json:"Message"`
	Time         string `json:"Time"`
}

type ApiResponse struct {
	Vendor string `json:"Vendor"`
}

type MapMessage struct {
	message     string
	routing_key string
	time        string
	valid       bool
}

type Device struct {
	Device_name string `json:"name"`
	Mac string `json:"mac"`
	Ip_address string `json:"address"`
	Alive bool `json:"alive"`
	Allowed int `json:"allowed"`
}

type DeviceFoundTopic struct {
	Device_name string `json:"name"`
	Ip_address string `json:"address"`
	Status int `json:"status"`
}

type DeviceUpdate struct {
	Name string `json:"name"`
	Mac string `json:"mac"`
	Status int `json:"status"`
	State string `json:"state"`
}

type DeviceRequest struct {
	Request_id uint32 `json:"id"`
	Name string `json:"name"`
	Mac string `json:"mac"`
}

type DeviceResponse struct {
	Request_id uint32 `json:"id"`
	Name string `json:"name"`
	Mac string `json:"mac"`
	Status string `json:"status"`
}

//Topics
const REQUESTDATA string = "Request.Data"
const AUTHENTICATIONREQUEST string = "Authentication.Request"
const DATAINFO string = "Data.Info"
const REQUESTACCESS string = "Request.Access"
const DEVICEADD string = "Device.Add"
const DEVICERESPONSE string = "Device.Response"
//
const FAILURENETWORK string = "Failure.Network"
const EVENTNAC string = "Event.NAC"
const REQUESTDATABASE string = "Request.Database"
const DATARESPONSE string = "Data.Response"
const DEVICEFOUND string = "Device.Found"
const AUTHENTICATIONRESPONSE string = "Authentication.Response"
const ACCESSRESPONSE string = "Access.Response"
const UNAUTHORISEDCONNECTION string = "Unauthorised.Connection"
const DEVICEREQUEST string = "Device.Request"
const DEVICEUPDATE string = "Device.Update"
//
const ACCESSFAIL string = "FAIL"
const ACCESSPASS string = "PASS"
const EXCHANGENAME string = "topics"
const EXCHANGETYPE string = "topic"
const TIMEFORMAT string = "20060102150405"
const COMPONENT string = "NAC"
const FAILUREPUBLISH string = "Failed to publish"
const SERVERERROR string = "Server is failing to send"
//
const ALLOWED int = 1
const BLOCKED int = 2
const DISCOVERED int = 3
const UNKNOWN int = 4
const ALLOWED_STRING string = "ALLOWED"
const BLOCKED_STRING string = "BLOCKED"
const UNKNOWN_STRING string = "UNKNOWN"
//
const START_ADDRESS string = "192.168.0."
//
var SubscribedMessagesMap map[uint32]*MapMessage
var DevicesList map[uint32]*Device
var key_id uint32 = 0
var device_id uint32 = 0
