// rabbitmq_test.go

package main

import (
	"strings"
	"testing"
)

// Check that State is set
// then run this test will prove it is set
func TestPublishFailRabbit(t *testing.T) {
	failure := "cheese"
	failure = messageFailure(true)
	if failure != "" {
		if strings.Contains(FAILUREPUBLISH, failure) {
			t.Error("Failure")
		} else if strings.Contains(FAILURECONVERT, failure) {
			t.Error("Failure")
		}
	}
}

func TestLogicNetwork(t *testing.T) {
	value := "{ 'time': 12:00:34, 'type': 'Camera', 'severity': 3 }"
	messages(FAILURENETWORK, value)
	checkState()
	if SubscribedMessagesMap[0].valid == true {
		t.Error("Failure")
	} else if SubscribedMessagesMap[0].routing_key != FAILURENETWORK {
		t.Log(SubscribedMessagesMap[0].routing_key)
		t.Error("Failure")
	}
}

func TestLogicNotExpected(t *testing.T) {
	value := "{ 'time': 12:00:34, 'type': 'Camera', 'severity': 3 }"
	messages("Event.DBM", value)
	checkState()
	if SubscribedMessagesMap[1].valid == true {
		t.Error("Failure")
	} else if SubscribedMessagesMap[1].routing_key == EVENTFH {
		t.Log(SubscribedMessagesMap[1].routing_key)
		t.Error("Failure")
	}
}

func TestLogicValid(t *testing.T) {
	value := "{ 'time': 12:00:34, 'type': 'Camera', 'severity': 3 }"
	messages("Event.DBM", value)
	if SubscribedMessagesMap[2].valid == false {
		t.Error("Failure")
	} else if SubscribedMessagesMap[2].routing_key == EVENTNAC {
		t.Log(SubscribedMessagesMap[2].routing_key)
		t.Error("Failure")
	}
}


func TestGetTime(t *testing.T) {
	time := getTime()
	if !strings.Contains(time, "2020") {
		t.Error("Failure")
	}
}

func TestGetTimeFail(t *testing.T) {
	time := getTime()
	if strings.Contains(time, ":") {
		t.Error("Failure")
	}
}

func TestEventNAC(t *testing.T) {
	valid := PublishEventNAC (COMPONENT, UPDATESTATEERROR, getTime())
	if valid != "" {
		t.Error("Failure")
	}
}


