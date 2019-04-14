package iu880b

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"./crc"
	"./slip"
	"./wimod"
)

func TestMain(t *testing.T) {
	performTest()
}

func TestMain2(t *testing.T) {
	performTest2()
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

func TestInherit(t *testing.T) {
	pingReq := wimod.NewPingReq()
	id := pingReq.GetID()
	fmt.Println(id)
}

func TestTime(t *testing.T) {
	now := time.Now().UTC()
	rtc := wimod.EncodeRTCTime(now)
	fmt.Printf("%s\n", now)
	fmt.Printf("%d\n", rtc)
	now2 := wimod.DecodeRTCTime(rtc)
	fmt.Printf("%s", now2)
}
