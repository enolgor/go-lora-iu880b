package iu880b

import (
	"fmt"

	"./crc"
	"./slip"
)

type HCIPacket struct {
	Payload []byte
	Dst     byte
	ID      byte
	Err     error
	CRC     []byte
}

func (p HCIPacket) String() string {
	return fmt.Sprintf("HCI: Dst=[%X] ID=[%X] Payload=[%X] CRC=[%X]", p.Dst, p.ID, p.Payload, p.CRC)
}

func HCIDecoder(slipPacketChan <-chan slip.SlipPacket, hciPacketChan chan<- HCIPacket) {
	for packet := range slipPacketChan {
		if packet.Err != nil {
			fmt.Printf("Packet Error: %s\n", packet.Err.Error())
		}
		bytes := packet.Buffer.Bytes()
		hciPacket := HCIPacket{}
		fmt.Printf("Packet Received: %X - ", bytes)
		if crc.CheckCRC16(bytes) {
			hciPacket.Dst = bytes[0]
			hciPacket.ID = bytes[1]
			hciPacket.Payload = bytes[2 : len(bytes)-2]
			hciPacket.CRC = bytes[len(bytes)-2:]
		} else {
			hciPacket.Err = fmt.Errorf("CRC check fail")
		}
		hciPacketChan <- hciPacket
	}
}
