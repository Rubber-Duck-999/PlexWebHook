package main

import (
	"os"

	"github.com/akamensky/argparse"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Warn("FH - Beginning to run Network Access Controller Program")
	parser := argparse.NewParser("file", "Config file for runtime purpose")
	// Create string flag
	f := parser.String("f", "config", &argparse.Options{Required: true, Help: "Necessary config"})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		log.Error(parser.Usage(err))
		os.Exit(2)
	}

	file := *f
	var data ConfigTypes
	if Exists(file) {
		GetData(&data, file)
		SetCodes(data.Settings.Code,
			data.Settings.Default_Pin)
		DevicesList = make(map[uint32]*DeviceFound)
		device_id = 0
		primary := DeviceFound{data.Primary.Name, 
			data.Primary.Mac,
			data.Primary.Ip,
			false}
		DevicesList[key_id] = &primary
		device_id++
		secondary := DeviceFound{data.Secondary.Name, 
			data.Secondary.Mac,
			data.Secondary.Ip,
			false}
		DevicesList[1] = &secondary
		device_id++
		tertiary := DeviceFound{data.Tertiary.Name, 
			data.Tertiary.Mac,
			data.Tertiary.Ip,
			false}
		DevicesList[2] = &tertiary
		device_id++
		checkDevices()
		Subscribe()
	} else {
		log.Error("File doesn't exist")
		os.Exit(2)
	}
	log.Trace(data.Settings.Code)
}
