package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/enolgor/wimod-lorawan-endnode-controller/controller"
	"github.com/enolgor/wimod-lorawan-endnode-controller/wimod"
	"github.com/tarm/serial"
)

// controller info -network -firmware -device ...
// controller join -type otaa|abp -appkey asdf -nwkskey asdf -appskey asdf
// controller send -enc ascii|hex|b64 -type u|c asdfasdf

const usageMessage = `
Usage: loractl <command> [<args>]

Available commands:
  info  Display information about the network/device
  join  Join a LoRa network
  send  Send a packet to the network
`

var infoCommand = flag.NewFlagSet("info", flag.ExitOnError)
var joinCommand = flag.NewFlagSet("join", flag.ExitOnError)
var sendCommand = flag.NewFlagSet("send", flag.ExitOnError)
var deactivateCommand = flag.NewFlagSet("deactivate", flag.ExitOnError)

var serialPort string

const (
	serialPortFlag  = "port"
	serialPortUsage = "Set serial port of LoRa EndNode device"
)

var infoNetwork bool

const (
	infoNetworkFlag        = "network"
	defaultInfoNetworkFlag = false
	infoNetworkUsage       = "Display lora network information"
)

var infoFirmware bool

const (
	infoFirmwareFlag        = "firmware"
	defaultInfoFirmwareFlag = false
	infoFirmwareUsage       = "Display firmware information"
)

var infoDevice bool

const (
	infoDeviceFlag        = "device"
	defaultInfoDeviceFlag = false
	infoDeviceUsage       = "Display device information"
)

var infoStatus bool

const (
	infoStatusFlag        = "status"
	defaultInfoStatusFlag = false
	infoStatusUsage       = "Display status information"
)

func init() {
	infoCommand.StringVar(&serialPort, serialPortFlag, "", serialPortUsage)
	infoCommand.BoolVar(&infoNetwork, infoNetworkFlag, defaultInfoNetworkFlag, infoNetworkUsage)
	infoCommand.BoolVar(&infoFirmware, infoFirmwareFlag, defaultInfoFirmwareFlag, infoFirmwareUsage)
	infoCommand.BoolVar(&infoDevice, infoDeviceFlag, defaultInfoDeviceFlag, infoDeviceUsage)
	infoCommand.BoolVar(&infoStatus, infoStatusFlag, defaultInfoStatusFlag, infoStatusUsage)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprint(os.Stderr, usageMessage)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "info":
		infoCommand.Parse(os.Args[2:])
		runInfoCommand()
	case "join":
	case "send":
	case "deactivate":
	default:
		fmt.Fprintf(os.Stderr, "%q is not a valid command\n", os.Args[1])
		fmt.Fprint(os.Stderr, usageMessage)
		os.Exit(1)
	}
}

func printDefaults(flagSet *flag.FlagSet) {
	fmt.Fprintf(os.Stderr, "\n%s command usage:\n\n", flagSet.Name())
	flagSet.PrintDefaults()
}

func checkSerialPort(flagSet *flag.FlagSet) {
	if serialPort == "" {
		fmt.Fprintln(os.Stderr, "Serial post must be specified")
		printDefaults(flagSet)
		os.Exit(1)
	}
}

func printErrorAndExit(e error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", e.Error())
	os.Exit(1)
}

func getController() *controller.WiModController {
	c := &serial.Config{Name: serialPort, Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		printErrorAndExit(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	return controller.NewController(config)
}

func dataRateString(idx byte) string {
	switch idx {
	case 0:
		return "LoRa SF12, 125kHz, 250bps"
	case 1:
		return "LoRa SF11, 125kHz, 440bps"
	case 2:
		return "LoRa SF10, 125kHz, 980bps"
	case 3:
		return "LoRa SF9, 125kHz, 1760bps"
	case 4:
		return "LoRa SF8, 125kHz, 3125bps"
	case 5:
		return "LoRa SF7, 125kHz, 5470bps"
	case 6:
		return "LoRa SF7, 250kHz, 11000bps"
	case 7:
		return "FSK 50K, -, 50000bps"
	}
	return "Unknown Data Rate"
}

func runInfoCommand() {
	checkSerialPort(infoCommand)
	if !infoNetwork && !infoFirmware && !infoDevice && !infoStatus {
		printDefaults(infoCommand)
		os.Exit(1)
	}
	controller := getController()

	if infoNetwork {
		req := wimod.NewGetNwkStatusReq()
		resp := wimod.NewGetNwkStatusResp()
		err := controller.Request(req, resp)
		if err != nil {
			printErrorAndExit(err)
		}
		var statusStr string
		switch resp.NetworkStatus {
		case wimod.LORAWAN_NETWORK_STATUS_INACTIVE:
			statusStr = "INACTIVE"
		case wimod.LORAWAN_NETWORK_STATUS_ACTIVE_ABP:
			statusStr = "ACTIVE (ABP)"
		case wimod.LORAWAN_NETWORK_STATUS_ACTIVE_OTAA:
			statusStr = "ACTIVE (OTAA)"
		case wimod.LORAWAN_NETWORK_STATUS_JOINING_OTAA:
			statusStr = "JOINING (OTAA)"
		}
		fmt.Fprint(os.Stdout, "\nNETWORK INFO:\n\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "Network Status:\t%s\n", statusStr)
		fmt.Fprintf(w, "Address:\t%08X\n", resp.Address)
		fmt.Fprintf(w, "Data Rate Index:\t%d - %s\n", resp.DataRateIdx, dataRateString(resp.DataRateIdx))
		fmt.Fprintf(w, "Power Level:\t%d dBm\n", resp.PowerLevel)
		fmt.Fprintf(w, "Max Payload Size:\t%d bytes\n", resp.MaxPayloadSize)
		w.Flush()
	}

	if infoFirmware {
		req := wimod.NewGetFWInfoReq()
		resp := wimod.NewGetFWInfoResp()
		err := controller.Request(req, resp)
		if err != nil {
			printErrorAndExit(err)
		}
		fmt.Fprint(os.Stdout, "\nFIRMWARE INFO:\n\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "Version:\t%d.%d.%d\n", resp.MajorVersion, resp.MinorVersion, resp.Build)
		fmt.Fprintf(w, "Build Date:\t%s\n", resp.BuildDate)
		fmt.Fprintf(w, "Firmware Image:\t%s\n", resp.FirmwareImage)
		w.Flush()
	}

	if infoDevice {
		req := wimod.NewGetDeviceInfoReq()
		resp := wimod.NewGetDeviceInfoResp()
		err := controller.Request(req, resp)
		if err != nil {
			printErrorAndExit(err)
		}
		var moduleTypeStr string
		switch resp.ModuleType {
		case 0x90:
			moduleTypeStr = "iM880A"
		case 0x92:
			moduleTypeStr = "iM880A-L"
		case 0x93:
			moduleTypeStr = "iU880A"
		case 0x98:
			moduleTypeStr = "iM880B-L"
		case 0x99:
			moduleTypeStr = "iU880B"
		case 0x9A:
			moduleTypeStr = "iM980A"
		case 0xA0:
			moduleTypeStr = "iM881A"
		}
		fmt.Fprint(os.Stdout, "\nDEVICE INFO:\n\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "Module Type:\t%s\n", moduleTypeStr)
		fmt.Fprintf(w, "Device Address:\t%08X\n", resp.DeviceAddress)
		fmt.Fprintf(w, "Device ID:\t%08X\n", resp.DeviceID)
		w.Flush()
	}

	/*if infoStatus {
		req := wimod.NewGetDeviceStatusReq()
		resp := wimod.NewGetDeviceStatusResp()
		err := controller.Request(req, resp)
		if err != nil {
			printErrorAndExit(err)
		}
		fmt.Fprint(os.Stdout, "\nSTATUS INFO:\n\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, ":\t%s\n", moduleTypeStr)
		fmt.Fprintf(w, "Device Address:\t%08X\n", resp.DeviceAddress)
		fmt.Fprintf(w, "Device ID:\t%08X\n", resp.DeviceID)
		w.Flush()
	}*/

}
