package wimod

import (
	"fmt"

	"../hci"
)

type WiModMessage interface {
	ID() byte
	Dst() byte
	Code() uint16
}

type wimodMessageImpl struct {
	code uint16
}

func (w *wimodMessageImpl) ID() byte {
	return byte(w.code & 0x00FF)
}

func (w *wimodMessageImpl) Dst() byte {
	return byte(w.code >> 8)
}

func (w *wimodMessageImpl) Code() uint16 {
	return w.code
}

type WiModMessageReq interface {
	WiModMessage
	Encode() ([]byte, error)
}

type WiModMessageResp interface {
	WiModMessage
	Decode(bytes []byte) error
}

type wimodMessageIndImpl struct {
	wimodMessageImpl
	status byte
}

func (w *wimodMessageIndImpl) Status() byte {
	return w.status
}

type WiModMessageInd interface {
	WiModMessage
	Decode(bytes []byte) error
	Status() byte
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
	status := hci.Payload[0]
	err := statusCheck(hci.Dst, status)
	if err != nil {
		return err
	}
	resp.Decode(hci.Payload[1:])
	return nil
}

func DecodeInd(hci *hci.HCIPacket) (WiModMessageInd, error) {
	code := (uint16(hci.Dst) << 8) + uint16(hci.ID)
	if !IsAlarm(code) {
		return nil, fmt.Errorf("Packet is not an event")
	}
	status := hci.Payload[0]
	err := statusCheck(hci.Dst, status)
	if err != nil {
		return nil, err
	}
	ind := alarmConstructors[code]()
	ind.Decode(hci.Payload) //INCLUDE STATUS
	return ind, nil
}

func statusCheck(dst, status byte) error {
	switch dst {
	case DEVMGMT_ID:
		switch status {
		case DEVMGMT_STATUS_OK:
			return nil
		case DEVMGMT_STATUS_ERROR:
			return fmt.Errorf("DEVMGMT_STATUS_ERROR")
		case DEVMGMT_STATUS_CMD_NOT_SUPPORTED:
			return fmt.Errorf("DEVMGMT_STATUS_ERROR")
		case DEVMGMT_STATUS_WRONG_PARAMETER:
			return fmt.Errorf("DEVMGMT_STATUS_ERROR")
		}
	case LORAWAN_ID:
		switch status {
		case LORAWAN_STATUS_OK:
			return nil
		case LORAWAN_STATUS_ERROR:
			return fmt.Errorf("LORAWAN_STATUS_ERROR")
		case LORAWAN_STATUS_CMD_NOT_SUPPORTED:
			return fmt.Errorf("LORAWAN_STATUS_CMD_NOT_SUPPORTED")
		case LORAWAN_STATUS_WRONG_PARAMETER:
			return fmt.Errorf("LORAWAN_STATUS_WRONG_PARAMETER")
		case LORAWAN_STATUS_WRONG_DEVICE_MODE:
			return fmt.Errorf("LORAWAN_STATUS_WRONG_DEVICE_MODE")
		case LORAWAN_STATUS_DEVICE_NOT_ACTIVATED:
			return fmt.Errorf("LORAWAN_STATUS_DEVICE_NOT_ACTIVATED")
		case LORAWAN_STATUS_DEVICE_BUSY:
			return fmt.Errorf("LORAWAN_STATUS_DEVICE_BUSY")
		case LORAWAN_STATUS_QUEUE_FULL:
			return fmt.Errorf("LORAWAN_STATUS_QUEUE_FULL")
		case LORAWAN_STATUS_LENGTH_ERROR:
			return fmt.Errorf("LORAWAN_STATUS_LENGTH_ERROR")
		case LORAWAN_STATUS_NO_FACTORY_SETTINGS:
			return fmt.Errorf("LORAWAN_STATUS_NO_FACTORY_SETTINGS")
		case LORAWAN_STATUS_CHANNEL_BLOCKED:
			return fmt.Errorf("LORAWAN_STATUS_CHANNEL_BLOCKED")
		case LORAWAN_STATUS_CHANNEL_NOT_AVAILABLE:
			return fmt.Errorf("LORAWAN_STATUS_CHANNEL_NOT_AVAILABLE")
		}
	}
	return fmt.Errorf("Unknown DST")
}
