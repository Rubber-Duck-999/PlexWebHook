// logicoptions_test.go

package main

import (
	"encoding/json"
	"testing"
)

func TestSwitchStatus(t *testing.T) {
	status := "BLOCKED"
	if convertStatus(status) != 2 {
		t.Error("Failure")
	}

	status = "GRANTED"
	if convertStatus(status) != 3 {
		t.Error("Failure")
	}
}

func TestStatusCheckMessageFalse(t *testing.T) {
	failure, _ := json.Marshal(&FailureNetwork{
		Time:         "",
		Failure_type: ""})
	s := string(failure)
	message := MapMessage{
		message:     s,
		routing_key: FAILURENETWORK,
		time:        getTime(),
		valid:       true,
	}
	if convertStatusMessage(message) {
		t.Error("Failure")
	}
}

func TestStatusCheckMessageTrue(t *testing.T) {
	status, _ := json.Marshal(&StatusFH{
		DailyFaults:  0,
		CommonFaults: ""})
	s := string(status)
	message := MapMessage{
		message:     s,
		routing_key: STATUSFH,
		time:        getTime(),
		valid:       true,
	}
	if !convertStatusMessage(message) {
		t.Error("Failure")
	}
}
