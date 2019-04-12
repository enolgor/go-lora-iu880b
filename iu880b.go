package iu880b

import (
	"bytes"
	"fmt"
	"log"

	"./crc"
	"./slip"
	"github.com/tarm/serial"
)

func performTest() {
	c := &serial.Config{Name: "COM3", Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	data := []byte{DEVMGMT_ID, DEVMGMT_MSG_GET_DEVICE_STATUS_REQ}
	crcval := crc.CalcCRC16(data)
	var buff bytes.Buffer
	data = append(data, byte(crcval&0xFF), byte((crcval>>8)&0xFF))
	buff.WriteByte(slip.SLIP_END)
	buff.Write(data)
	buff.WriteByte(slip.SLIP_END)
	fmt.Printf("%x\n", buff.Bytes())
	sendWakeUp(s)
	_, err = s.Write(buff.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	slipPacketChan := make(chan slip.SlipPacket)
	hciPacketChan := make(chan HCIPacket)
	go slip.SlipDecoder(s, slipPacketChan)
	go HCIDecoder(slipPacketChan, hciPacketChan)
	for packet := range hciPacketChan {
		switch packet.Dst {
		case 0x01:
			fmt.Println("DVMGMT message")
			switch packet.ID {
			case DEVMGMT_MSG_GET_DEVICE_STATUS_RSP:
				deviceStatus := &WiModDevStatus{}
				deviceStatus.Decode(packet.Payload)
				fmt.Println(deviceStatus)
				break
			}
			break
		case 0x10:
			fmt.Println("LoraWan HCI Message")
			break
		default:
			fmt.Printf("Unrecognized Endpoint Identifier byte: %X\n", packet.Dst)
		}
		fmt.Println(packet)
	}
}

func sendWakeUp(s *serial.Port) {
	for i := 0; i < 40; i++ {
		s.Write([]byte{slip.SLIP_END})
	}
}
