package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var _port string
var _device_status string
var _messages_done bool

type allData []DataInfo

var data_messages = allData{
}

func init() {
	_port = "0"
}

func SetPort(port string) {
	log.Debug("Port set")
	_port = port
}

func isValidGUID(guid string) bool {
	return true
}

func device_add(w http.ResponseWriter, r *http.Request) {
	var device DeviceAdd
	req_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		json.Unmarshal(req_body, &device)
		if isValidGUID(device.GUID) {
			log.Debug("Received Device Name: ", device.Name)
			PublishDeviceUpdate(device.Name, device.Mac,
								device.Status, r.Method)
			w.WriteHeader(http.StatusOK)
		} else {
			log.Error("Invalid GUID")
		}
	}
}

func user_add(w http.ResponseWriter, r *http.Request) {
	var user UserAdd
	req_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		json.Unmarshal(req_body, &user)
		if isValidGUID(user.GUID) {
			log.Debug("Received User Name: ", user.User)
			//PublishUserUpdate(user.User, user.Pin, r.Method)
			w.WriteHeader(http.StatusOK)
		} else {
			log.Error("Invalid GUID")
		}
	}
}

func data_request(w http.ResponseWriter, r *http.Request) {
	log.Warn("Received data message")
	var request RequestData
	req_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		json.Unmarshal(req_body, &request)
		if isValidGUID(request.GUID) {
			PublishRequestDatabase(request.Request_id, request.Time_from, 
							request.Time_to, request.Type)
			time.Sleep(2 * time.Second)
			if _messages_done && data_messages != nil {
				json.NewEncoder(w).Encode(data_messages)
			} else {
				time.Sleep(2 * time.Second)
				if _messages_done && data_messages != nil {
					json.NewEncoder(w).Encode(data_messages)
				}
			}

		} else {
			log.Error("Invalid GUID")
		}
	}
}

func http_server() {
	router := mux.NewRouter().StrictSlash(true)
	// Set up of methods
	router.HandleFunc("/device", device_add).Methods("POST")
	router.HandleFunc("/device", device_add).Methods("PATCH")
	router.HandleFunc("/device", device_add).Methods("DELETE")
	//
	router.HandleFunc("/user", user_add).Methods("POST")
	router.HandleFunc("/user", user_add).Methods("PATCH")
	router.HandleFunc("/user", user_add).Methods("DELETE")	
	//
	router.HandleFunc("/data", data_request).Methods("GET")
	//	
	log.Fatal(http.ListenAndServe(":" + _port, router))
}