package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var _endpoint string

func setEndpoint(endpoint string) {
	_endpoint = endpoint
}

func apiCall(q string, name string) {
	call_allowed := true
	client := http.Client{}
	if call_allowed {
		req, _ := http.NewRequest("POST", _endpoint+name, strings.NewReader(q))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		resp, err := client.Do(req)
		if err != nil {
			log.Warn("Error updating server")
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warn("Body reading error")
		}
		timeString := time.Now().String()
		PublishFailureNetwork(timeString, name)
	}
}

func postAlarmEvent(event AlarmEvent) {
	q := url.Values{}
	q.Add("user", "User")
	q.Add("state", "OFF")
	name := "/alarmEvent"
	apiCall(q.Encode(), name)
}

func postAccess() {
	q := url.Values{}
	q.Add("access_granted", _statusUP.LastAccessGranted)
	q.Add("access_denied", _statusUP.LastAccessBlocked)
	q.Add("last_user", _statusUP.LastUser)
	name := "/access"
	apiCall(q.Encode(), name)
}

func postDevice() {
	q := url.Values{}
	q.Add("allowed", strconv.Itoa(_statusNAC.DailyAllowedDevices))
	q.Add("blocked", strconv.Itoa(_statusNAC.DailyBlockedDevices))
	q.Add("unknown", strconv.Itoa(_statusNAC.DailyUnknownDevices))
	name := "/device"
	apiCall(q.Encode(), name)
}

func postFault() {
	q := url.Values{}
	q.Add("name", _statusFH.LastFault)
	name := "/fault"
	apiCall(q.Encode(), name)
}

func postHardware() {
	q := url.Values{}
	q.Add("cpu_temp", strconv.Itoa(_statusSYP.Temperature))
	q.Add("cpu_usage", strconv.Itoa(_statusSYP.HighestUsage))
	q.Add("memory", strconv.Itoa(_statusSYP.MemoryLeft))
	name := "/hardware"
	apiCall(q.Encode(), name)
}

func postMotion() {
	q := url.Values{}
	name := "/motion"
	apiCall(q.Encode(), name)
}
