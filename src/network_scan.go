package main

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// All ye beware, there be dragons below...

import (
    "context"
    "io/ioutil"
    "net/http"
    "time"
	"os/exec"
    "strings"

    log "github.com/sirupsen/logrus"
    "github.com/Ullaakut/nmap"
)

const (
    ProtocolICMP = 1
)

// Default to listen on all IPv4 interfaces
var ListenAddr = "0.0.0.0"
var _statusNAC StatusNAC

func init() {
    _statusNAC = StatusNAC{
		DevicesActive: 0,
		DailyBlockedDevices: 0,
		DailyUnknownDevices: 0,
		DailyAllowedDevices: 0,
		TimeEscConnected: "N/A"}
}

func runARP() {
    log.Debug("### Running ARP ###")
    data, err := exec.Command("arp", "-a").Output()
	if err != nil {
		log.Error(err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// strip brackets around IP
		ip := strings.Replace(fields[1], "(", "", -1)
        ip = strings.Replace(ip, ")", "", -1)
        new_device := true
        for id := range DevicesList {
            if DevicesList[id].Ip_address == ip {
                log.Trace("Device found in Arp table")
                DevicesList[id].Alive = true
                new_device = false
            }
        }
        mac := fields[3]
        if new_device == true {
            if mac != "<incomplete>" {
                log.Debug("Adding device ip: ", ip)
                response, err := http.Get("https://api.macvendors.com/" + mac)
                if err != nil {
                    log.Error("The HTTP request failed with error \n", err)
                } else {
                    data, _ := ioutil.ReadAll(response.Body)
                    log.Trace(response)
                    log.Debug("Vendor Name: ", string(data))
                    device := Device{string(data), mac, ip, true, DISCOVERED}
                    DevicesList[device_id] = &device
                    device_id++
                }
                time.Sleep(1 * time.Second)
            }
        }
        
    }
}

func nmap_scan() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    // Equivalent to `/usr/local/bin/nmap -p 80,443,843 google.com facebook.com youtube.com`,
    // with a 2 minute timeout.
    scanner, err := nmap.NewScanner(
        nmap.WithTargets("192.168.0.0-255"),
        nmap.WithPorts("80,443,843"),
        nmap.WithContext(ctx),
    )
    if err != nil {
        log.Error("Unable to create nmap scanner: ", err)
    }

    result, warnings, err := scanner.Run()
    if err != nil {
        log.Error("Unable to run nmap scan: ", err)
    }

    if warnings != nil {
        log.Error("Warnings: ", warnings)
    }

    // Use the results to print an example output
    for _, host := range result.Hosts {
        if len(host.Ports) == 0 || len(host.Addresses) == 0 {
            continue
        }

        log.Debug("Host: ", host.Addresses[0])

        for _, port := range host.Ports {
            if port.State.String() != "closed" {
                log.Debug("# Port ID: ", port.ID)
                log.Debug("# Protocol: ", port.Protocol)
                log.Debug("# State: ", port.State)
                log.Debug("# Service: ", port.Service.Name)
            }
        }
    }

    log.Debug("Nmap done: ", len(result.Hosts), " hosts up scanned in seconds ",  result.Stats.Finished.Elapsed)
}

func checkDevices() {
    done := false
    for {
        if done == false {
            nmap_scan()
            runARP()
            log.Warn("### Devices ###")
            _statusNAC.DevicesActive = 0
            for id := range DevicesList {
                log.Warn("Device - ", DevicesList[id].Device_name, " : ",
                    DevicesList[id].Ip_address, " : ", 
                    DevicesList[id].Mac, " : ",
                    DevicesList[id].Alive, " : ",
                    DevicesList[id].Allowed)
                if DevicesList[id].Alive == true {
                    if DevicesList[id].Allowed == BLOCKED {
                        PublishDeviceFound(DevicesList[id].Device_name,
                            DevicesList[id].Ip_address,
                            DevicesList[id].Allowed)
                    } else if DevicesList[id].Allowed == DISCOVERED {
                        PublishDeviceRequest(id,
                            DevicesList[id].Device_name,
                            DevicesList[id].Mac)
                    } else if DevicesList[id].Allowed == UNKNOWN {
                        PublishDeviceFound(DevicesList[id].Device_name,
                            DevicesList[id].Ip_address,
                            DevicesList[id].Allowed)
                    } else {
                        log.Debug("Device is Allowed so moving to next")
                    }
                    _statusNAC.DevicesActive++
                }
            }
            log.Debug("### End of ARP ###")
            log.Debug("Starting Status NAC publish")
            log.Debug("Current message : ", _statusNAC)
            log.Debug("### End of Status ###")
            PublishStatusNAC()
            done = true
        } else {
            done = false
        }
        time.Sleep(2 * time.Minute)
    }
}