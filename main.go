package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/enolgor/wimod-lorawan-endnode-controller/controller"
	"github.com/enolgor/wimod-lorawan-endnode-controller/wimod"
	"github.com/tarm/serial"
)

// controller info -network -firmware -device -radio
// controller join -type otaa|abp -appkey asdf -nwkskey asdf -appskey asdf -eui asdf
// controller send -enc ascii|hex|b64 -type u|c asdfasdf -port 1
// controller synctime
// controller deactivate

const usageMessage = `
Usage: loractl <command> [<args>]

Available commands:
  info        Display information about the network/device
  join        Join a LoRa network
  send        Send a packet to the network
  synctime    Synchronize time with the host machine
  deactivate  Deactivate device
`

var infoCommand = flag.NewFlagSet("info", flag.ExitOnError)
var joinCommand = flag.NewFlagSet("join", flag.ExitOnError)
var sendCommand = flag.NewFlagSet("send", flag.ExitOnError)
var synctimeCommand = flag.NewFlagSet("synctime", flag.ExitOnError)
var deactivateCommand = flag.NewFlagSet("deactivate", flag.ExitOnError)

var serialPort string

const (
	serialPortFlag        = "serialport"
	defaultSerialPortFlag = ""
	serialPortUsage       = "Set serial port of LoRa EndNode device"
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

var infoRadio bool

const (
	infoRadioFlag        = "radio"
	defaultInfoRadioFlag = false
	infoRadioUsage       = "Display radio information"
)

var joinType string

const (
	joinTypeFlag        = "type"
	defaultJoinTypeFlag = ""
	joinTypeUsage       = "Specify join type: abp|otaa"
)

var appKey string

const (
	appKeyFlag        = "appkey"
	defaultAppKeyFlag = ""
	appKeyUsage       = "Specify APP Key (required for otaa)"
)

var appEUI string

const (
	appEUIFlag        = "appeui"
	defaultAppEUIFlag = ""
	appEUIUsage       = "Specify APP EUI (required for otaa)"
)

var address string

const (
	addressFlag        = "address"
	defaultAddressFlag = ""
	addressFlagUsage   = "Specify address of device (required for abp)"
)

var appSessKey string

const (
	appSessKeyFlag        = "appsesskey"
	defaultAppSessKeyFlag = ""
	appSessKeyUsage       = "Specify APP Session Key (required for abp)"
)

var nwkSessKey string

const (
	nwkSessKeyFlag        = "nwksesskey"
	defaultNwkSessKeyFlag = ""
	nwkSessKeyUsage       = "Specify Network Session Key (required for abp)"
)

var sendEnc string

const (
	sendEncFlag        = "enc"
	defaultSendEncFlag = ""
	sendEncUsage       = "Specify encoding of payload: ascii|hex|b64"
)

var sendType string

const (
	sendTypeFlag        = "type"
	defaultSendTypeFlag = ""
	sendTypeUsage       = "Specify type of packet sent: c|u"
)

var sendPayload string

const (
	sendPayloadFlag        = "payload"
	defaultSendPayloadFlag = ""
	sendPayloadUsage       = "Payload to send"
)

var sendPort uint

const (
	sendPortFlag        = "loraport"
	defaultSendPortFlag = 0
	sendPortUsage       = "Specify port: 0-255"
)

func init() {
	infoCommand.StringVar(&serialPort, serialPortFlag, defaultSerialPortFlag, serialPortUsage)
	infoCommand.BoolVar(&infoNetwork, infoNetworkFlag, defaultInfoNetworkFlag, infoNetworkUsage)
	infoCommand.BoolVar(&infoFirmware, infoFirmwareFlag, defaultInfoFirmwareFlag, infoFirmwareUsage)
	infoCommand.BoolVar(&infoDevice, infoDeviceFlag, defaultInfoDeviceFlag, infoDeviceUsage)
	infoCommand.BoolVar(&infoStatus, infoStatusFlag, defaultInfoStatusFlag, infoStatusUsage)
	infoCommand.BoolVar(&infoRadio, infoRadioFlag, defaultInfoRadioFlag, infoRadioUsage)

	joinCommand.StringVar(&serialPort, serialPortFlag, defaultSerialPortFlag, serialPortUsage)
	joinCommand.StringVar(&joinType, joinTypeFlag, defaultJoinTypeFlag, joinTypeUsage)
	joinCommand.StringVar(&appKey, appKeyFlag, defaultAppKeyFlag, appKeyUsage)
	joinCommand.StringVar(&appEUI, appEUIFlag, defaultAppEUIFlag, appEUIUsage)
	joinCommand.StringVar(&address, addressFlag, defaultAddressFlag, addressFlagUsage)
	joinCommand.StringVar(&appSessKey, appSessKeyFlag, defaultAppSessKeyFlag, appSessKeyUsage)
	joinCommand.StringVar(&nwkSessKey, nwkSessKeyFlag, defaultNwkSessKeyFlag, nwkSessKeyUsage)

	sendCommand.StringVar(&serialPort, serialPortFlag, defaultSerialPortFlag, serialPortUsage)
	sendCommand.StringVar(&sendEnc, sendEncFlag, defaultSendEncFlag, sendEncUsage)
	sendCommand.StringVar(&sendType, sendTypeFlag, defaultSendTypeFlag, sendTypeUsage)
	sendCommand.StringVar(&sendPayload, sendPayloadFlag, defaultSendPayloadFlag, sendPayloadUsage)
	sendCommand.UintVar(&sendPort, sendPortFlag, defaultSendPortFlag, sendPortUsage)

	deactivateCommand.StringVar(&serialPort, serialPortFlag, defaultSerialPortFlag, serialPortUsage)

	synctimeCommand.StringVar(&serialPort, serialPortFlag, defaultSerialPortFlag, serialPortUsage)

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
		joinCommand.Parse(os.Args[2:])
		runJoinCommand()
	case "send":
		sendCommand.Parse(os.Args[2:])
		runSendCommand()
	case "synctime":
		synctimeCommand.Parse(os.Args[2:])
		runSynctimeCommand()
	case "deactivate":
		deactivateCommand.Parse(os.Args[2:])
		runDeactivateCommand()
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

func getTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
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
	if !infoNetwork && !infoFirmware && !infoDevice && !infoStatus && !infoRadio {
		printDefaults(infoCommand)
		os.Exit(1)
	}
	controller := getController()
	w := getTabWriter()
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
		fmt.Fprint(w, "\nNETWORK INFO:\n\n")
		fmt.Fprintf(w, "Network Status:\t%s\n", statusStr)
		fmt.Fprintf(w, "Address:\t%08X\n", resp.Address)
		fmt.Fprintf(w, "Data Rate:\t%d - %s\n", resp.DataRateIdx, dataRateString(resp.DataRateIdx))
		fmt.Fprintf(w, "Power Level:\t%d dBm\n", resp.PowerLevel)
		fmt.Fprintf(w, "Max Payload Size:\t%d bytes\n", resp.MaxPayloadSize)
	}

	if infoFirmware {
		req := wimod.NewGetFWInfoReq()
		resp := wimod.NewGetFWInfoResp()
		err := controller.Request(req, resp)
		if err != nil {
			printErrorAndExit(err)
		}
		fmt.Fprint(w, "\nFIRMWARE INFO:\n\n")
		fmt.Fprintf(w, "Version:\t%d.%d.%d\n", resp.MajorVersion, resp.MinorVersion, resp.Build)
		fmt.Fprintf(w, "Build Date:\t%s\n", resp.BuildDate)
		fmt.Fprintf(w, "Firmware Image:\t%s\n", resp.FirmwareImage)
	}

	if infoDevice {
		req := wimod.NewGetDeviceInfoReq()
		resp := wimod.NewGetDeviceInfoResp()
		err := controller.Request(req, resp)
		if err != nil {
			printErrorAndExit(err)
		}
		euiReq := wimod.NewGetDeviceEUIReq()
		euiResp := wimod.NewGetDeviceEUIResp()
		err = controller.Request(euiReq, euiResp)
		if err != nil {
			printErrorAndExit(err)
		}
		opModeReq := wimod.NewGetOPModeReq()
		opModeResp := wimod.NewGetOPModeResp()
		err = controller.Request(opModeReq, opModeResp)
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
		default:
			moduleTypeStr = "Unknown"
		}
		var opModeStr string
		switch opModeResp.Mode {
		case 0x00:
			opModeStr = "Default/Standard"
		case 0x01:
			opModeStr = "Reserved"
		case 0x02:
			opModeStr = "Reserved"
		case 0x03:
			opModeStr = "Customer"
		default:
			opModeStr = "Unknown"
		}
		fmt.Fprint(w, "\nDEVICE INFO:\n\n")
		fmt.Fprintf(w, "Module Type:\t%s\n", moduleTypeStr)
		fmt.Fprintf(w, "Device Address:\t%08X\n", resp.DeviceAddress)
		fmt.Fprintf(w, "Device ID:\t%08X\n", resp.DeviceID)
		fmt.Fprintf(w, "Device EUI:\t%v\n", euiResp.EUI)
		fmt.Fprintf(w, "Operation Mode:\t%s\n", opModeStr)
	}

	if infoStatus {
		req := wimod.NewGetDeviceStatusReq()
		resp := wimod.NewGetDeviceStatusResp()
		err := controller.Request(req, resp)
		if err != nil {
			printErrorAndExit(err)
		}
		fmt.Fprint(w, "\nSTATUS INFO:\n\n")
		fmt.Fprintf(w, "System Tick Resolution:\t%d ms\n", resp.SystemTickResolution)
		fmt.Fprintf(w, "System Ticks:\t%d\n", resp.SystemTicks)
		fmt.Fprintf(w, "Target Time:\t%s\n", resp.TargetTime)
		sysstatus := "ok"
		if resp.NVMStatus&0x01 == 1 {
			sysstatus = "error"
		}
		radiostatus := "ok"
		if (resp.NVMStatus>>1)&0x01 == 1 {
			radiostatus = "error"
		}
		fmt.Fprintf(w, "NVM Status:\tsystem:%s, radio:%s\n", sysstatus, radiostatus)
		fmt.Fprintf(w, "Battery Level:\t%d mV\n", resp.BatteryLevel)
		fmt.Fprintf(w, "Extra Status:\t0x%04X\n", resp.ExtraStatus)
		fmt.Fprintf(w, "Tx U-Data:\t%d packets\n", resp.TxUData)
		fmt.Fprintf(w, "Tx C-Data:\t%d packets\n", resp.TxCData)
		fmt.Fprintf(w, "Tx Error:\t%d packets\n", resp.TxError)
		fmt.Fprintf(w, "Rx1 U-Data:\t%d packets\n", resp.Rx1UData)
		fmt.Fprintf(w, "Rx1 C-Data:\t%d packets\n", resp.Rx1CData)
		fmt.Fprintf(w, "Rx1 MIC Error:\t%d packets\n", resp.Rx1MICError)
		fmt.Fprintf(w, "Rx2 U-Data:\t%d packets\n", resp.Rx2UData)
		fmt.Fprintf(w, "Rx2 C-Data:\t%d packets\n", resp.Rx2CData)
		fmt.Fprintf(w, "Rx2 MIC Error:\t%d packets\n", resp.Rx2MICError)
		fmt.Fprintf(w, "Tx Join:\t%d packets\n", resp.TxJoin)
		fmt.Fprintf(w, "Rx Accept:\t%d packets\n", resp.RxAccept)
	}

	if infoRadio {
		req := wimod.NewGetRStackConfigReq()
		resp := wimod.NewGetRStackConfigResp()
		err := controller.Request(req, resp)
		if err != nil {
			printErrorAndExit(err)
		}
		enabled := func(e bool) string {
			if e {
				return "enabled"
			}
			return "disabled"
		}
		fmt.Fprint(w, "\nRADIO INFO:\n\n")
		fmt.Fprintf(w, "Default Data Rate:\t%d - %s\n", resp.DefaultDataRateIdx, dataRateString(resp.DefaultDataRateIdx))
		fmt.Fprintf(w, "TX Power Level:\t%d dBm\n", resp.TXPowerLevel)
		fmt.Fprintf(w, "Adaptative Data Rate:\t%s\n", enabled(resp.AdaptativeDataRate))
		fmt.Fprintf(w, "Duty Cycle Control:\t%s\n", enabled(resp.DutyCycleControl))
		fmt.Fprintf(w, "Class C:\t%s\n", enabled(resp.ClassC))
		fmt.Fprintf(w, "MAC Events:\t%s\n", enabled(resp.MACEvents))
		fmt.Fprintf(w, "Extended HCI:\t%s\n", enabled(resp.ExtendedHCI))
		fmt.Fprintf(w, "Automatic Power Saving:\t%s\n", enabled(resp.AutomaticPowerSaving))
		fmt.Fprintf(w, "Max Retransmissions:\t%d\n", resp.MaxRetransmissions)
		fmt.Fprintf(w, "Band Index:\t%d\n", resp.BandIdx)
		fmt.Fprintf(w, "Header MAC Cmd Capacity:\t%d\n", resp.HeaderMACCmdCapacity)
	}

	w.Flush()
}

func runDeactivateCommand() {
	checkSerialPort(deactivateCommand)
	controller := getController()
	deactivateReq := wimod.NewDeactivateDeviceReq()
	deactivateResp := wimod.NewDeactivateDeviceResp()
	err := controller.Request(deactivateReq, deactivateResp)
	if err != nil {
		printErrorAndExit(err)
	}
	w := getTabWriter()
	fmt.Fprintf(w, "Successfully deactivated\n")
	w.Flush()
}

func runSynctimeCommand() {
	checkSerialPort(synctimeCommand)
	controller := getController()
	reqSet := wimod.NewSetRTCReq(time.Now().UTC())
	respSet := wimod.NewSetRTCResp()
	err := controller.Request(reqSet, respSet)
	if err != nil {
		printErrorAndExit(err)
	}
	reqGet := wimod.NewGetRTCReq()
	respGet := wimod.NewGetRTCResp()
	err = controller.Request(reqGet, respGet)
	if err != nil {
		printErrorAndExit(err)
	}
	w := getTabWriter()
	fmt.Fprintf(w, "Time synced:\t%s\n", respGet.Time)
	w.Flush()
}

func runJoinCommand() {
	checkSerialPort(joinCommand)
	switch joinType {
	case "abp":
		if appSessKey == "" || nwkSessKey == "" || address == "" {
			fmt.Fprintln(os.Stderr, "For ABP join type, appsesskey, nwksesskey and address must be specified")
			printDefaults(joinCommand)
			os.Exit(1)
		}
		err := abpJoin()
		if err != nil {
			printErrorAndExit(err)
		}
	case "otaa":
		if appKey == "" || appEUI == "" {
			fmt.Fprintln(os.Stderr, "For OTAA join type, appkey and appeui must be specified")
			printDefaults(joinCommand)
			os.Exit(1)
		}
		err := otaaJoin()
		if err != nil {
			printErrorAndExit(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "Join type must be abp or otaa")
		printDefaults(joinCommand)
		os.Exit(1)
	}
}

func otaaJoin() error {
	controller := getController()
	eui, err := wimod.ParseEUI(appEUI)
	if err != nil {
		return err
	}
	key, err := wimod.ParseKey(appKey)
	if err != nil {
		return err
	}
	nwkStatusReq := wimod.NewGetNwkStatusReq()
	nwkStatusResp := wimod.NewGetNwkStatusResp()
	err = controller.Request(nwkStatusReq, nwkStatusResp)
	if err != nil {
		return err
	}
	if nwkStatusResp.NetworkStatus != wimod.LORAWAN_NETWORK_STATUS_INACTIVE {
		printErrorAndExit(fmt.Errorf("device is already joined or joining, deactivate first"))
	}
	joinParamReq := wimod.NewSetJoinParamReq(eui, key)
	joinParamResp := wimod.NewSetJoinParamResp()
	err = controller.Request(joinParamReq, joinParamResp)
	if err != nil {
		return err
	}
	joinReq := wimod.NewJoinNetworkReq()
	joinResp := wimod.NewJoinNetworkResp()
	err = controller.Request(joinReq, joinResp)
	if err != nil {
		return err
	}
	joinTxEvent := wimod.NewJoinNetworkTxInd()
	joinedEvent := wimod.NewJoinNetworkInd()
	err = controller.ReadSpecificInd(joinTxEvent)
	if err != nil {
		return err
	}
	w := getTabWriter()
	fmt.Fprintf(w, "Join packet successfully sent\n")
	err = controller.ReadSpecificInd(joinedEvent)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Device successfully joined\n")
	fmt.Fprintf(w, "Address:\t%08X\n", joinedEvent.Address)
	w.Flush()
	return nil
}

func abpJoin() error {
	controller := getController()
	keyApp, err := wimod.ParseKey(appSessKey)
	if err != nil {
		return err
	}
	keyNwk, err := wimod.ParseKey(nwkSessKey)
	if err != nil {
		return err
	}
	addr, err := strconv.ParseUint(address, 16, 32)
	if err != nil {
		return err
	}
	nwkStatusReq := wimod.NewGetNwkStatusReq()
	nwkStatusResp := wimod.NewGetNwkStatusResp()
	err = controller.Request(nwkStatusReq, nwkStatusResp)
	if err != nil {
		return err
	}
	if nwkStatusResp.NetworkStatus != wimod.LORAWAN_NETWORK_STATUS_INACTIVE {
		printErrorAndExit(fmt.Errorf("device is already joined or joining, deactivate first"))
	}
	activateReq := wimod.NewActivateDeviceReq(uint32(addr), keyApp, keyNwk)
	activateResp := wimod.NewActivateDeviceResp()
	err = controller.Request(activateReq, activateResp)
	if err != nil {
		return err
	}
	w := getTabWriter()
	fmt.Fprintf(w, "Device successfully activated\n")
	w.Flush()
	return nil
}

func runSendCommand() {
	checkSerialPort(sendCommand)
	if sendType != "c" && sendType != "u" {
		printErrorAndExit(fmt.Errorf("send type should be (c)onfirmed or (u)nconfirmed"))
	}
	if sendEnc != "ascii" && sendEnc != "hex" && sendEnc != "b64" {
		printErrorAndExit(fmt.Errorf("encoding should be ascii, hex or b64"))
	}
	if sendPayload == "" {
		printErrorAndExit(fmt.Errorf("payload is empty"))
	}
	if sendPort > 255 {
		printErrorAndExit(fmt.Errorf("port should be from 0 to 255"))
	}
	port := byte(sendPort)
	var payload []byte
	var err error
	switch sendEnc {
	case "ascii":
		payload = []byte(sendPayload)
	case "hex":
		payload, err = hex.DecodeString(sendPayload)
		if err != nil {
			printErrorAndExit(err)
		}
	case "b64":
		payload, err = base64.StdEncoding.DecodeString(sendPayload)
		if err != nil {
			printErrorAndExit(err)
		}
	}
	switch sendType {
	case "c":
		err = sendConfirmed(port, payload)
		if err != nil {
			printErrorAndExit(err)
		}
	case "u":
		err = sendUnconfirmed(port, payload)
		if err != nil {
			printErrorAndExit(err)
		}
	}
}

func sendUnconfirmed(port byte, payload []byte) error {
	controller := getController()
	udataReq := wimod.NewSendUDataReq(port, payload)
	udataResp := wimod.NewSendUDataResp()
	err := controller.Request(udataReq, udataResp)
	if err != nil {
		return err
	}
	udataTxInd := wimod.NewSendUDataTxInd()
	err = controller.ReadSpecificInd(udataTxInd)
	if err != nil {
		return err
	}
	w := getTabWriter()
	fmt.Fprintf(w, "Unconfirmed data successfully sent\n")
	w.Flush()
	return nil
}

func sendConfirmed(port byte, payload []byte) error {
	return nil
}
