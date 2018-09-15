package main

import (
	"errors"
	"fmt"
	"github.com/adelolmo/hd-idle/hdidle"
	"github.com/jasonlvhit/gocron"
	"os"
	"strconv"
)

const (
	DEFAULT_IDLE_TIME = 600
)

func main() {

	if os.Getenv("START_HD_IDLE") == "false" {
		println("START_HD_IDLE=false exiting now.")
		os.Exit(0)
	}

	defaultConf := hdidle.DefaultConf{
		Idle:        DEFAULT_IDLE_TIME,
		CommandType: hdidle.SCSI,
		Debug:       false,
	}
	var config = &hdidle.Config{
		Devices:  []hdidle.DeviceConf{},
		Defaults: defaultConf,
	}
	var deviceConf *hdidle.DeviceConf

	for index, arg := range os.Args[1:] {
		switch arg {
		case "-a":
			if deviceConf != nil {
				config.Devices = append(config.Devices, *deviceConf)
			}

			name := os.Args[index+2]
			deviceConf = &hdidle.DeviceConf{
				Name:        name,
				Idle:        config.Defaults.Idle,
				CommandType: config.Defaults.CommandType,
			}
			break

		case "-i":
			s := os.Args[index+2]
			idle, err := strconv.Atoi(s)
			if err != nil {
				println(errors.New(fmt.Sprintf("Wrong idle_time -i %d. Must be a number", idle)))
				os.Exit(1)
			}
			if deviceConf == nil {
				config.Defaults.Idle = idle
				break
			}
			deviceConf.Idle = idle
			break

		case "-c":
			command := os.Args[index+2]
			switch command {
			case hdidle.SCSI:
			case hdidle.ATA:
				if deviceConf == nil {
					config.Defaults.CommandType = command
					break
				}
				deviceConf.CommandType = command
				break
			default:
				println(errors.New(fmt.Sprintf("Wrong command_type -c %s. Must be one of: scsi, ata", command)))
				os.Exit(1)
			}
			break

		case "-l":
			config.Defaults.LogFile = os.Args[index+2]
			break

		case "-d":
			config.Defaults.Debug = true
			break

		case "h":
			println("usage: hd-idle [-t <disk>] [-a <name>] [-i <idle_time>] [-c <command_type>] [-l <logfile>] [-d] [-h]\n")
			os.Exit(0)
		}
	}

	if deviceConf != nil {
		config.Devices = append(config.Devices, *deviceConf)
	}
	println(config.String())

	gocron.Every(60).Seconds().Do(hdidle.ObserveDiskActivity, config)
	gocron.NextRun()
	<-gocron.Start()
}
