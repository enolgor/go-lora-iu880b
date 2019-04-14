package wimod

import (
	"fmt"

	"../hci"
)

type WiModMessage interface {
	GetID() byte
	GetDst() byte
}

type wimodmessage struct {
	id  byte
	dst byte
}

func (w *wimodmessage) GetID() byte {
	return w.id
}

func (w *wimodmessage) GetDst() byte {
	return w.dst
}

type WiModMessageReq interface {
	WiModMessage
	Encode() ([]byte, error)
}

type WiModMessageResp interface {
	WiModMessage
	Decode(bytes []byte) error
}

func EncodeReq(req WiModMessageReq) (*hci.HCIPacket, error) {
	payload, err := req.Encode()
	if err != nil {
		return nil, err
	}
	return &hci.HCIPacket{Dst: req.GetDst(), ID: req.GetID(), Payload: payload}, nil
}

func DecodeResp(hci *hci.HCIPacket, resp WiModMessageResp) error {
	if hci.Dst != resp.GetDst() || hci.ID != resp.GetID() {
		return fmt.Errorf("Wrong DST or ID")
	}
	status := hci.Payload[0]
	if status != DEVMGMT_STATUS_OK {
		switch status {
		case DEVMGMT_STATUS_ERROR:
			return fmt.Errorf("DEVMGMT_STATUS_ERROR")
		case DEVMGMT_STATUS_CMD_NOT_SUPPORTED:
			return fmt.Errorf("DEVMGMT_STATUS_CMD_NOT_SUPPORTED")
		case DEVMGMT_STATUS_WRONG_PARAMETER:
			return fmt.Errorf("DEVMGMT_STATUS_WRONG_PARAMETER")
		default:
			return fmt.Errorf("Unknown")
		}
	}
	resp.Decode(hci.Payload[1:])
	return nil
}
