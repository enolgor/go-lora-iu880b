package wimod

import (
	"fmt"

	"github.com/enolgor/wimod-lorawan-endnode-controller/hci"
)

type WiModMessage interface {
	Code() uint16
	ID() byte
	Dst() byte
	Init()
}

type wimodMessageImpl struct {
	code uint16
}

func (w *wimodMessageImpl) Code() uint16 {
	return w.code
}

func (w *wimodMessageImpl) ID() byte {
	return byte(w.code & 0x00FF)
}

func (w *wimodMessageImpl) Dst() byte {
	return byte(w.code >> 8)
}

type wimodMessageStatusImpl struct {
	wimodMessageImpl
	Status byte
}

type WiModMessageReq interface {
	WiModMessage
	Encode() ([]byte, error)
}

type WiModMessageResp interface {
	WiModMessage
	Decode(bytes []byte) error
}

type WiModMessageInd interface {
	WiModMessage
	Decode(bytes []byte) error
}

func EncodeReq(req WiModMessageReq) (*hci.HCIPacket, error) {
	payload, err := req.Encode()
	if err != nil {
		return nil, err
	}
	return &hci.HCIPacket{Dst: req.Dst(), ID: req.ID(), Payload: payload}, nil
}

func DecodeResp(hci *hci.HCIPacket, resp WiModMessageResp) error {
	if hci.Dst != resp.Dst() || hci.ID != resp.ID() {
		return fmt.Errorf("Wrong DST or ID")
	}
	return resp.Decode(hci.Payload)
}

func DecodeInd(hci *hci.HCIPacket) (WiModMessageInd, error) {
	code := (uint16(hci.Dst) << 8) + uint16(hci.ID)
	if !IsAlarm(code) {
		return nil, fmt.Errorf("Packet is not an event")
	}
	ind := alarmConstructors[code]()
	err := ind.Decode(hci.Payload) //INCLUDE STATUS
	return ind, err
}

func DecodeSpecificInd(hci *hci.HCIPacket, ind WiModMessageInd) error {
	if hci.Dst != ind.Dst() || hci.ID != ind.ID() {
		return fmt.Errorf("Wrong DST or ID")
	}
	return ind.Decode(hci.Payload)
}
