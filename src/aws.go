package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var _endpoint string

func setEndpoint(endpoint string) {
	_endpoint = endpoint
}

func apiCall(req *http.Request) bool {
	call_allowed := true
	client := http.Client{}
	if call_allowed {
		resp, err := client.Do(req)
		if err != nil {
			log.Warn("Error updating server")
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warn("Body reading error")
		}
		return true
	}
	return false
}

func driveUpdateStatus() {
	t := time.Now()
 
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	if h == 2 && m == 0 {
		if s == 0 && s < 5 {
			log.Debug("Posting daily status")
			postDailyStatus()
		}
	}
	if m == 1 {
		log.Debug("Posting status")
		postStatus()
		time.Sleep(2 * time.Minute)
	}
}

func postAlarmEvent(event AlarmEvent) {
	q := url.Values{}
	q.Add("user", "User")
	q.Add("state", "OFF")
	req, _ := http.NewRequest("POST", _endpoint+"/alarmEvent", strings.NewReader(q.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if apiCall(req) {
		log.Debug("Request Successful")
	} else {
		log.Error("Request Failed on POST AlarmEvent")
		time_string := time.Now().String()
		PublishFailureNetwork(time_string, "AlarmEvent")
	}
}

func postStatus() {
	q := url.Values{}
	t := time.Now()
	q.Add("created_date", t.Format("2006-01-02"))
	q.Add("motion_detected", "")
	q.Add("access_granted", "")
	q.Add("access_denied", _statusUP.LastAccessBlocked)
	q.Add("last_fault", "N/A")
	q.Add("last_user", _statusUP.LastUser)
	q.Add("cpu_temp", strconv.Itoa(_statusSYP.Temperature))
	q.Add("cpu_usage", strconv.Itoa(_statusSYP.HighestUsage))
	q.Add("memory", strconv.Itoa(_statusSYP.MemoryLeft))
	req, _ := http.NewRequest("POST", _endpoint+"/status", strings.NewReader(q.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if apiCall(req) {
		log.Debug("Request Successful")
	} else {
		log.Error("Request Failed on POST Status")
		time_string := time.Now().String()
		PublishFailureNetwork(time_string, "Status")
	}
}

func postDailyStatus() {
	q := url.Values{}
	t := time.Now()
	q.Add("created_date", t.Format("2006-01-02"))
	q.Add("allowed", strconv.Itoa(_statusNAC.DailyAllowedDevices))
	q.Add("blocked", strconv.Itoa(_statusNAC.DailyBlockedDevices))
	q.Add("unknown", strconv.Itoa(_statusNAC.DailyUnknownDevices))
	q.Add("total_events", "")
	q.Add("common_event", "")
	q.Add("total_faults", strconv.Itoa(_statusFH.DailyFaults))
	q.Add("common_fault", _statusFH.CommonFaults)
	req, _ := http.NewRequest("POST", _endpoint+"/dailyStatus", strings.NewReader(q.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if apiCall(req) {
		log.Debug("Request Successful")
	} else {
		log.Error("Request Failed on POST DailyStatus")
		time_string := time.Now().String()
		PublishFailureNetwork(time_string, "Daily Status")
	}
}
