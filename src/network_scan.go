package main

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// All ye beware, there be dragons below...

import (
    "time"
    "os"
    "fmt"
	"net"
	"os/exec"
	"strings"
    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"

    log "github.com/sirupsen/logrus"
)

const (
    ProtocolICMP = 1
    network_scan_minute = 15
)

// Default to listen on all IPv4 interfaces
var ListenAddr = "0.0.0.0"

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
                log.Debug("Cannot find MAC going to ping: ", ip)
                dst, dur, err := Ping(ip)
                log.Trace("Ping ", dst, " : ", dur, " : ", err)
            } else {
                log.Debug("Adding device: ", ip)
                device := Device{"New", mac, ip, true, DISCOVERED}
                DevicesList[device_id] = &device
                device_id++
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
    reply := make([]byte, 1500)
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
    for {
        t := time.Now()
        min := t.Minute()
        done := false
        mod := min % network_scan_minute 
        log.Trace("Minute is: ", min)
        if mod == 0 && !done {
            for id := range DevicesList {
                dst, dur, err := Ping(DevicesList[id].Ip_address)
                if err != nil {
                    log.Error("Error in Ping on device: ", 
                        DevicesList[id].Device_name)
                    DevicesList[id].Alive = false
                } else {
                    log.Trace("Ping ", DevicesList[id].Device_name, " : ", 
                        dst, " : ", dur)
                    DevicesList[id].Alive = true
                }
            }
            runARP()
            runARP()
            log.Warn("### Devices ###")
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
                            DevicesList[id].Ip_address,
                            DevicesList[id].Mac)
                    } else if DevicesList[id].Allowed == UNKNOWN {
                        PublishDeviceFound(DevicesList[id].Device_name,
                            DevicesList[id].Ip_address,
                            DevicesList[id].Allowed)
                    } else {
                        log.Debug("Device is Allowed so moving to next")
                    }
                }
            }
            done = true
        } else if mod == 0 && done {
            log.Debug("Not the right time to scan")
        } else {
            done = false
        }
    }
}