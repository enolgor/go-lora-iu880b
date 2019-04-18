package main

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/enolgor/wimod-lorawan-endnode-controller/controller"
	"github.com/enolgor/wimod-lorawan-endnode-controller/crc"
	"github.com/enolgor/wimod-lorawan-endnode-controller/slip"
	"github.com/enolgor/wimod-lorawan-endnode-controller/wimod"
	"github.com/tarm/serial"
)

func serialConfig() *serial.Config {
	return &serial.Config{Name: "COM3", Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
}

func TestEUI(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	euiReq := wimod.NewGetDeviceEUIReq()
	euiResp := wimod.NewGetDeviceEUIResp()
	err = controller.Request(euiReq, euiResp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(euiResp)
}

func TestSendUData(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	udataReq := wimod.NewSendUDataReq(2, []byte("Hola Mundo!"))
	udataResp := wimod.NewSendUDataResp()
	err = controller.Request(udataReq, udataResp)
	if err != nil {
		fmt.Println(udataResp)
		log.Fatal(err)
	}
	udataTxInd := wimod.NewSendUDataTxInd()
	err = controller.ReadSpecificInd(udataTxInd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(udataTxInd)
}

func TestJoin(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	EUI := wimod.DecodeEUI([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22})
	Key := wimod.DecodeKey([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	joinParamReq := wimod.NewSetJoinParamReq(EUI, Key)
	joinParamResp := wimod.NewSetJoinParamResp()
	err = controller.Request(joinParamReq, joinParamResp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Join Params Set")
	joinReq := wimod.NewJoinNetworkReq()
	joinResp := wimod.NewJoinNetworkResp()
	err = controller.Request(joinReq, joinResp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Join Command Sent")
	joinTxEvent := wimod.NewJoinNetworkTxInd()
	joinedEvent := wimod.NewJoinNetworkInd()
	err = controller.ReadSpecificInd(joinTxEvent)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(joinTxEvent)
	err = controller.ReadSpecificInd(joinedEvent)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(joinedEvent)
}

func TestNetworkStatus(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	nwkStatusReq := wimod.NewGetNwkStatusReq()
	nwkStatusResp := wimod.NewGetNwkStatusResp()
	err = controller.Request(nwkStatusReq, nwkStatusResp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(nwkStatusResp)
}
func TestReset(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	go func() {
		ind, err := controller.ReadInd()
		if err != nil {
			fmt.Printf("Event Error: %s\n", err.Error())
		} else {
			fmt.Printf("Received event: %v\n", ind)
		}
	}()
	resetReq := wimod.NewResetReq()
	resetResp := wimod.NewResetResp()
	err = controller.Request(resetReq, resetResp)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)
}
func TestEvents(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	setAlarmFromNow(controller, 2*time.Second)
	for i := 0; i < 10; i++ {
		go func(id int) {
			ind, err := controller.ReadInd()
			if err != nil {
				fmt.Printf("Event Error: %s\n", err.Error())
			} else {
				fmt.Printf("Received from routine %d: %v\n", id, ind)
			}
		}(i)
	}
	time.Sleep(5 * time.Second)
	controller.Close()
}

func TestOpMode(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	requestAndPrintOpMode(controller)
	req := wimod.NewSetOPModeReq(wimod.DEVMGMT_OPMODE_STANDARD)
	resp := wimod.NewSetOPModeResp()
	err = controller.Request(req, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	time.Sleep(2 * time.Second)
	requestAndPrintOpMode(controller)
	controller.Close()
}

func requestAndPrintOpMode(controller *controller.WiModController) {
	req := wimod.NewGetOPModeReq()
	resp := wimod.NewGetOPModeResp()
	err := controller.Request(req, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

func TestGetTime(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	req := wimod.NewGetRTCReq()
	resp := wimod.NewGetRTCResp()
	err = controller.Request(req, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	controller.Close()
}

func TestSetTime(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	req := wimod.NewSetRTCReq(time.Now().UTC())
	resp := wimod.NewSetRTCResp()
	err = controller.Request(req, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	controller.Close()
}

func TestAlarm(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	go func() {
		alarm := wimod.NewRTCAlarmInd()
		err := controller.ReadSpecificInd(alarm) // print events
		if err != nil {
			fmt.Printf("Event Error: %s\n", err.Error())
		} else {
			fmt.Println(alarm)
		}
	}()
	clearAlarm(controller)
	requestAndPrintAlarm(controller) // 0s

	setAlarmFromNow(controller, 5*time.Second) // alarm in 5s

	time.Sleep(2 * time.Second)

	requestAndPrintAlarm(controller) // 2s, alarm is on
	clearAlarm(controller)
	requestAndPrintAlarm(controller) // 2s, alarm is off

	setAlarmFromNow(controller, 3*time.Second) // alarm in 3s
	requestAndPrintAlarm(controller)           // 2s, alarm is on

	time.Sleep(5 * time.Second)

	// 7s, alarm shouldve been printed

	controller.Close()
}

func requestAndPrintAlarm(controller *controller.WiModController) {
	reqGetAlarm := wimod.NewGetRTCAlarmReq()
	respGetAlarm := wimod.NewGetRTCAlarmResp()
	err := controller.Request(reqGetAlarm, respGetAlarm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(respGetAlarm)
}

func setAlarmFromNow(controller *controller.WiModController, d time.Duration) {
	now := time.Now().UTC().Add(d)
	reqSetAlarm := wimod.NewSetRTCAlarmReq(wimod.AlarmSingle, byte(now.Hour()), byte(now.Minute()), byte(now.Second()))
	respSetAlarm := wimod.NewSetRTCAlarmResp()
	err := controller.Request(reqSetAlarm, respSetAlarm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(respSetAlarm)
}

func clearAlarm(controller *controller.WiModController) {
	reqClearAlarm := wimod.NewClearRTCAlarmReq()
	respClearAlarm := wimod.NewClearRTCAlarmResp()
	err := controller.Request(reqClearAlarm, respClearAlarm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(respClearAlarm)
}

func TestBlocked(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{s, 1, false}
	controller := controller.NewController(config)
	now := time.Now().UTC().Add(2 * time.Second)
	req := wimod.NewSetRTCAlarmReq(wimod.AlarmSingle, byte(now.Hour()), byte(now.Minute()), byte(now.Second()))
	resp := wimod.NewSetRTCAlarmResp()
	err = controller.Request(req, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	time.Sleep(3 * time.Second)
	now = time.Now().UTC().Add(2 * time.Second)
	req = wimod.NewSetRTCAlarmReq(wimod.AlarmSingle, byte(now.Hour()), byte(now.Minute()), byte(now.Second()))
	resp = wimod.NewSetRTCAlarmResp()
	err = controller.Request(req, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	time.Sleep(6 * time.Second)
	req2 := wimod.NewGetRTCReq()
	resp2 := wimod.NewGetRTCResp()
	err = controller.Request(req2, resp2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp2)
	controller.Close()
}

func TestInfo(t *testing.T) {
	c := serialConfig()
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &controller.WiModControllerConfig{Stream: s}
	controller := controller.NewController(config)
	infoReq := wimod.NewGetDeviceInfoReq()
	infoResp := wimod.NewGetDeviceInfoResp()
	err = controller.Request(infoReq, infoResp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(infoResp)
}

func TestCRC(t *testing.T) {
	data := []byte{0x01, 0x01, 0x34, 0xFF, 0x56}
	crcval := crc.CalcCRC16(data)
	data = append(data, byte(crcval&0xFF), byte((crcval>>8)&0xFF))
	fmt.Printf("%X\n", data)
	res := crc.CheckCRC16(data)
	if !res {
		t.Error("CRC failed")
	}
}

func TestSLIP(t *testing.T) {
	data := []byte{0xaf, 0x2d, slip.SLIP_END, 0x3f, 0x11, slip.SLIP_ESC, 0xaa}
	buff := slip.SlipEncode(data)
	reader := bytes.NewReader(buff)
	slipDecoder := slip.NewDecoder(reader)
	pckt, err := slipDecoder.Read()
	if err != nil {
		t.Fatalf("Error decoding slip packet: %s", err.Error())
	}
	equals := bytes.Equal(pckt, data)
	if !equals {
		t.Log("Encoded byte slice does not equal decoded byte slice")
		t.FailNow()
	}
	t.Log("Byte slice encoded and decoded successfully")
}

func TestINFOSlip(t *testing.T) {
	infoReq := wimod.GetDeviceInfoReq{}
	hci, err := wimod.EncodeReq(&infoReq)
	if err != nil {
		t.Fatal(err)
	}
	slipPacket := slip.SlipEncode(hci.Encode())
	fmt.Printf("%X\n", slipPacket)
}

func TestTime(t *testing.T) {
	now := time.Now().UTC()
	rtc := wimod.EncodeRTCTime(now)
	fmt.Printf("%s\n", now)
	fmt.Printf("%d\n", rtc)
	now2 := wimod.DecodeRTCTime(rtc)
	fmt.Printf("%s", now2)
}

func TestKey(t *testing.T) {
	k := wimod.DecodeKey([]byte{0xAA, 0x23, 0x24, 0xFF, 0xAA, 0x23, 0x24, 0x5D, 0x44})
	fmt.Println(k)
	bytes := wimod.EncodeKey(&k)
	fmt.Printf("0x%02X\n", bytes)
	fmt.Println(wimod.DecodeKey(bytes))
}

func TestSizes(t *testing.T) {
	b := []byte{0x11, 0x22}
	fmt.Printf("%X\n", b[2:])
}
