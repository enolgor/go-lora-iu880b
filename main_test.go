package iu880b

import (
	"fmt"
	"testing"

	"./crc"
)

func TestMain(t *testing.T) {
	performTest()
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
