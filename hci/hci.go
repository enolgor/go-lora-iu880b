package hci

import (
	"bytes"
	"fmt"

	"../crc"
	"../slip"
)

type HCIPacket struct {
	Dst     byte
	ID      byte
	Payload []byte
}

func (p HCIPacket) String() string {
	return fmt.Sprintf("HCI: Dst=[%X] ID=[%X] Payload=[%X]", p.Dst, p.ID, p.Payload)
}

type HciDecoder struct {
	slipDecoder *slip.SlipDecoder
}

func NewDecoder(slipDecoder *slip.SlipDecoder) HciDecoder {
	return HciDecoder{slipDecoder}
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

func (h *HciDecoder) Read(hciPacket *HCIPacket) error {
	payload, err := h.slipDecoder.Read()
	if err != nil {
		return err
	}
	return decodeHCI(payload, hciPacket)
}

func decodeHCI(payload []byte, hciPacket *HCIPacket) error {
	if crc.CheckCRC16(payload) {
		hciPacket.Dst = payload[0]
		hciPacket.ID = payload[1]
		hciPacket.Payload = payload[2 : len(payload)-2]
		return nil
	}
	return fmt.Errorf("CRC check fail")
}
