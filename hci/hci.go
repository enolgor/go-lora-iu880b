package hci

import (
	"bytes"
	"fmt"

	"../crc"
)

type HCIPacket struct {
	Dst     byte
	ID      byte
	Payload []byte
}

func (p HCIPacket) String() string {
	return fmt.Sprintf("HCI: Dst=[%X] ID=[%X] Payload=[%X]", p.Dst, p.ID, p.Payload)
}

func (hci *HCIPacket) Encode() []byte {
	var buff bytes.Buffer
	buff.WriteByte(hci.Dst)
	buff.WriteByte(hci.ID)
	buff.Write(hci.Payload)
	crcval := crc.CalcCRC16(buff.Bytes())
	buff.WriteByte(byte(crcval & 0xFF))
	buff.WriteByte(byte((crcval >> 8) & 0xFF))
	return buff.Bytes()
}

func (hciPacket *HCIPacket) Decode(payload []byte) error {
	if crc.CheckCRC16(payload) {
		hciPacket.Dst = payload[0]
		hciPacket.ID = payload[1]
		hciPacket.Payload = payload[2 : len(payload)-2]
		return nil
	}
	return fmt.Errorf("CRC check fail")
}
