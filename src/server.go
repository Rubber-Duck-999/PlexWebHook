package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var _port string
var _device_status string

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
		}
	}
}

func http_server() {
	router := mux.NewRouter().StrictSlash(true)
	// Set up of methods
	router.HandleFunc("/device", device_add).Methods("POST")
	router.HandleFunc("/device", device_add).Methods("PATCH")
	router.HandleFunc("/device", device_add).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":" + _port, router))
}