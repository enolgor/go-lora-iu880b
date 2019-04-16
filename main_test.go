package iu880b

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"./crc"
	"./slip"
	"./wimod"
	"github.com/tarm/serial"
)

func TestOpMode(t *testing.T) {
	c := &serial.Config{Name: "COM3", Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &WiModControllerConfig{Stream: s}
	controller := NewController(config)
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

func requestAndPrintOpMode(controller *WiModController) {
	req := wimod.NewGetOPModeReq()
	resp := wimod.NewGetOPModeResp()
	err := controller.Request(req, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

func TestGetTime(t *testing.T) {
	c := &serial.Config{Name: "COM3", Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &WiModControllerConfig{Stream: s}
	controller := NewController(config)
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
	c := &serial.Config{Name: "COM3", Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &WiModControllerConfig{Stream: s}
	controller := NewController(config)
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
	c := &serial.Config{Name: "COM3", Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &WiModControllerConfig{Stream: s}
	controller := NewController(config)
	go func() {
		ind := controller.ReadInd() // print events
		fmt.Println(ind)
	}()

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

func requestAndPrintAlarm(controller *WiModController) {
	reqGetAlarm := wimod.NewGetRTCAlarmReq()
	respGetAlarm := wimod.NewGetRTCAlarmResp()
	err := controller.Request(reqGetAlarm, respGetAlarm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(respGetAlarm)
}

func setAlarmFromNow(controller *WiModController, d time.Duration) {
	now := time.Now().UTC().Add(d)
	reqSetAlarm := wimod.NewSetRTCAlarmReq(wimod.AlarmSingle, byte(now.Hour()), byte(now.Minute()), byte(now.Second()))
	respSetAlarm := wimod.NewSetRTCAlarmResp()
	err := controller.Request(reqSetAlarm, respSetAlarm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(respSetAlarm)
}

func clearAlarm(controller *WiModController) {
	reqClearAlarm := wimod.NewClearRTCAlarmReq()
	respClearAlarm := wimod.NewClearRTCAlarmResp()
	err := controller.Request(reqClearAlarm, respClearAlarm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(respClearAlarm)
}

func TestBlocked(t *testing.T) {
	c := &serial.Config{Name: "COM3", Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	config := &WiModControllerConfig{s, 1, false}
	controller := NewController(config)
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

func TestINFO(t *testing.T) {
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
