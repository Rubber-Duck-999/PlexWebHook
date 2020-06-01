package main

type ConfigTypes struct {
	Settings struct {
		Code    string `yaml:"code"`
		Port    string `yaml:"port"`
	} `yaml:"settings"`
}

// ESC Messages

type DeviceAdd struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
	Mac string `json:"mac"`
	Status int `json:"status"`
}

type UserAdd struct {
	GUID string `json:"guid"`
	User string `json:"user"`
	Pin int `json:"pin"`
}

type RequestData struct {
	GUID string `json:"guid"`
	Request_id int `json:"request_id"`
	Time_from string `json:"time_from"`
	Time_to string `json:"time_to"`
	EventTypeId string `json:"event_type_id"`
}

// End of ESC messages

type FailureNetwork struct {
	Time         string `json:"time"`
	Failure_type string `json:"type_of_failure"`
}

type RequestDatabase struct {
	Request_id int `json:"request_id"`
	Time_from string `json:"time_from"`
	Time_to string `json:"time_to"`
	EventTypeId string `json:"event_type_id"`
}

type DataInfo struct {
	Id int `json:"_id"`
	Message_num int `json:"_messageNum"`
	Total_message int `json:"_totalMessage"`
    Message string `json:"_topicMessage"`
	Time string `json:"_timeSent"`
}

type UnauthorisedConnection struct {
	Mac string `json:"mac"`
	Time string `json:"time"`
	Alive bool `json:"alive"`
}

type EventNAC struct {
	Component    string `json:"component"`
	Message      string `json:"message"`
	Time         string `json:"time"`
	EventTypeId  string `json:"event_type_id"`
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
	New bool `json:"new"`
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

type StatusNAC struct {
	DevicesActive int `json:"devices_active"`
	DailyBlockedDevices int `json:"blocked"`
	DailyUnknownDevices int `json:"unknown"`
	DailyAllowedDevices int `json:"allowed"`
	TimeEscConnected string `json:"time"`
}

//Topics
const REQUESTDATA string = "Request.Data"
const AUTHENTICATIONREQUEST string = "Authentication.Request"
const DATAINFO string = "Data.Info"
const DEVICEADD string = "Device.Add"
const DEVICERESPONSE string = "Device.Response"
//
const FAILURENETWORK string = "Failure.Network"
const EVENTNAC string = "Event.NAC"
const REQUESTDATABASE string = "Request.Database"
const DATARESPONSE string = "Data.Response"
const DEVICEFOUND string = "Device.Found"
const AUTHENTICATIONRESPONSE string = "Authentication.Response"
const UNAUTHORISEDCONNECTION string = "Unauthorised.Connection"
const DEVICEREQUEST string = "Device.Request"
const DEVICEUPDATE string = "Device.Update"
const STATUSNAC string = "Status.NAC"
//
const ACCESSFAIL string = "FAIL"
const ACCESSPASS string = "PASS"
const EXCHANGENAME string = "topics"
const EXCHANGETYPE string = "topic"
const TIMEFORMAT string = "20060102150405"
const COMPONENT string = "NAC"
const FAILUREPUBLISH string = "Failed to publish"
const UNKNOWN_DEVICE string = "New device connected - "
//
const ALLOWED int = 1
const BLOCKED int = 2
const DISCOVERED int = 3
const UNKNOWN int = 4
const ALLOWED_STRING string = "ALLOWED"
const BLOCKED_STRING string = "BLOCKED"
const DISCOVERED_STRING string = "DISCOVERED"
const UNKNOWN_STRING string = "UNKNOWN"
//
const START_ADDRESS string = "192.168.0."
//
var SubscribedMessagesMap map[uint32]*MapMessage
var DevicesList map[uint32]*Device
var key_id uint32 = 0
var device_id uint32 = 0
