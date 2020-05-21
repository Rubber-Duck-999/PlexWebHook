package main

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// All ye beware, there be dragons below...

import (
    "io/ioutil"
    "net/http"
    "time"
    "os"
    "fmt"
	"net"
	"os/exec"
    "strings"
    "strconv"
    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"

    log "github.com/sirupsen/logrus"
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
    log.Debug("Running ARP")
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
                log.Trace("Node found in Arp table")
                DevicesList[id].Alive = true
                new_device = false
            }
        }
        mac := fields[3]
        if new_device == true {
            if mac == "<incomplete>" {
                log.Trace("Cannot find MAC going to ping: ", ip)
                dst, dur, err := Ping(ip)
                log.Trace("Ping ", dst, " : ", dur, " : ", err)
            } else {
                log.Debug("Adding device: ", ip)
                response, err := http.Get("https://api.macvendors.com/" + mac)
                if err != nil {
                    log.Error("The HTTP request failed with error \n", err)
                } else {
                    data, _ := ioutil.ReadAll(response.Body)
                    log.Trace(response)
                    log.Debug("Device is actually: ", string(data))
                    device := Device{string(data), mac, ip, true, DISCOVERED}
                    DevicesList[device_id] = &device
                    device_id++
                }
                time.Sleep(1 * time.Second)
            }
        }
        
    }
}

func Ping(addr string) (*net.IPAddr, time.Duration, error) {
    // Start listening for icmp replies
    c, err := icmp.ListenPacket("ip4:icmp", ListenAddr)
    if err != nil {
        return nil, 0, err
    }
    defer c.Close()

    // Resolve any DNS (if used) and get the real IP of the target
    dst, err := net.ResolveIPAddr("ip4", addr)
    if err != nil {
        panic(err)
        return nil, 0, err
    }

    // Make a new ICMP message
    m := icmp.Message{
        Type: ipv4.ICMPTypeEcho, Code: 0,
        Body: &icmp.Echo{
            ID: os.Getpid() & 0xffff, Seq: 1,
            Data: []byte(""),
        },
    }
    b, err := m.Marshal(nil)
    if err != nil {
        return dst, 0, err
    }

    // Send it
    start := time.Now()
    n, err := c.WriteTo(b, dst)
    if err != nil {
        return dst, 0, err
    } else if n != len(b) {
        return dst, 0, fmt.Errorf("got %v; want %v", n, len(b))
    }

    // Wait for a reply
    reply := make([]byte, 3000)
    err = c.SetReadDeadline(time.Now().Add(10 * time.Second))
    if err != nil {
        return dst, 0, err
    }
    n, peer, err := c.ReadFrom(reply)
    if err != nil {
        return dst, 0, err
    }
    duration := time.Since(start)

    rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
    if err != nil {
        return dst, 0, err
    }
    switch rm.Type {
    case ipv4.ICMPTypeEchoReply:
        return dst, duration, nil
    default:
        return dst, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
    }
}

func checkDevices() {
    done := false
    for {
        if done == false {
            for addr := 0; addr < 32; addr++ {
                s := strconv.Itoa(addr)
                address:= START_ADDRESS + s
                dst, dur, err := Ping(address)
                if err != nil {
                    log.Trace("No reply")
                } else {
                    log.Debug("Ping ", dst, " : ", dur, " : ", err)
                }
            }
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
            log.Debug("### End ###")
            log.Debug("Starting Status NAC publish")
            log.Debug("Current message : ", _statusNAC)
            PublishStatusNAC()
            done = true
        } else {
            done = false
        }
        time.Sleep(2 * time.Minute)
    }
}