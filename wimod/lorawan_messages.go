package wimod

import (
	"encoding/binary"
	"fmt"
)

// LORAWAN_MSG_ACTIVATE_DEVICE_REQ
// LORAWAN_MSG_ACTIVATE_DEVICE_RSP
// LORAWAN_MSG_SET_JOIN_PARAM_REQ

type SetJoinParamReq struct {
	wimodMessageImpl
	AppEUI EUI
	AppKey Key
}

func NewSetJoinParamReq(appEUI EUI, appKey Key) *SetJoinParamReq {
	req := &SetJoinParamReq{}
	req.code = LORAWAN_MSG_SET_JOIN_PARAM_REQ
	req.AppEUI = appEUI
	req.AppKey = appKey
	return req
}

func (p *SetJoinParamReq) String() string {
	return fmt.Sprintf("SetJoinParamReq[AppEUI: %v, AppKey: %v]", p.AppEUI, p.AppKey)
}

func (p *SetJoinParamReq) Encode() ([]byte, error) {
	buff := EncodeEUI(&p.AppEUI)
	buff = append(buff, EncodeKey(&p.AppKey)...)
	return buff, nil
}

// LORAWAN_MSG_SET_JOIN_PARAM_RSP

type SetJoinParamResp struct {
	wimodMessageStatusImpl
}

func NewSetJoinParamResp() *SetJoinParamResp {
	resp := &SetJoinParamResp{}
	resp.code = LORAWAN_MSG_SET_JOIN_PARAM_RSP
	return resp
}

func (p *SetJoinParamResp) String() string {
	return fmt.Sprintf("SetJoinParamResp[]")
}

func (p *SetJoinParamResp) Decode(payload []byte) error {
	p.status = payload[0]
	return lorawanStatusCheck(p.status)
}

// LORAWAN_MSG_JOIN_NETWORK_REQ

type JoinNetworkReq struct {
	wimodMessageImpl
}

func NewJoinNetworkReq() *JoinNetworkReq {
	req := &JoinNetworkReq{}
	req.code = LORAWAN_MSG_JOIN_NETWORK_REQ
	return req
}

func (p *JoinNetworkReq) String() string {
	return fmt.Sprintf("JoinNetworkReq[]")
}

func (p *JoinNetworkReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// LORAWAN_MSG_JOIN_NETWORK_RSP

type JoinNetworkResp struct {
	wimodMessageStatusImpl
}

func NewJoinNetworkResp() *JoinNetworkResp {
	resp := &JoinNetworkResp{}
	resp.code = LORAWAN_MSG_JOIN_NETWORK_RSP
	return resp
}

func (p *JoinNetworkResp) String() string {
	return fmt.Sprintf("JoinNetworkResp[]")
}

func (p *JoinNetworkResp) Decode(payload []byte) error {
	p.status = payload[0]
	return lorawanStatusCheck(p.status)
}

// LORAWAN_MSG_JOIN_NETWORK_TX_IND

type JoinNetworkTxInd struct {
	wimodMessageStatusImpl
	ChannelIdx       byte
	DataRateIdx      byte
	NumTxPackets     byte
	TRXPowerLevel    byte
	RFMessageAirtime uint32
}

func NewJoinNetworkTxInd() *JoinNetworkTxInd {
	ind := &JoinNetworkTxInd{}
	ind.code = LORAWAN_MSG_JOIN_NETWORK_TX_IND
	return ind
}

func (p *JoinNetworkTxInd) String() string {
	return fmt.Sprintf("JoinNetworkTxInd[Status: 0x%02X, ChannelIdx: %d, DataRateIdx: %d, NumTxPackets: %d, TRXPowerLevel: %d, RFMessageAirtime: %d]", p.status, p.ChannelIdx, p.DataRateIdx, p.NumTxPackets, p.TRXPowerLevel, p.RFMessageAirtime)
}

func (p *JoinNetworkTxInd) Decode(bytes []byte) error {
	p.status = bytes[0]
	if p.status != LORAWAN_MSG_JOIN_NETWORK_TX_IND_STATUS_OK && p.status != LORAWAN_MSG_JOIN_NETWORK_TX_IND_STATUS_OK_ATTACHMENT {
		p.status = LORAWAN_MSG_JOIN_NETWORK_TX_IND_STATUS_ERROR
	}
	if p.status == LORAWAN_MSG_JOIN_NETWORK_TX_IND_STATUS_OK_ATTACHMENT {
		p.ChannelIdx = bytes[1]
		p.DataRateIdx = bytes[2]
		p.NumTxPackets = bytes[3]
		p.TRXPowerLevel = bytes[4]
		p.RFMessageAirtime = binary.LittleEndian.Uint32(bytes[5:9])
	}
	return nil
}

// LORAWAN_MSG_JOIN_NETWORK_IND

type JoinNetworkInd struct {
	wimodMessageStatusImpl
	Address     uint32
	ChannelIdx  byte
	DataRateIdx byte
	RSSI        byte
	SNR         byte
	RxSlot      byte
}

func NewJoinNetworkInd() *JoinNetworkInd {
	ind := &JoinNetworkInd{}
	ind.code = LORAWAN_MSG_JOIN_NETWORK_IND
	return ind
}

func (p *JoinNetworkInd) String() string {
	return fmt.Sprintf("JoinNetworkInd[Status: 0x%02X, Address: 0x%08X, ChannelIdx: %d, DataRateIdx: %d, RSSI: %d, SNR: %d, RxSlot: %d]", p.status, p.Address, p.ChannelIdx, p.DataRateIdx, p.RSSI, p.SNR, p.RxSlot)
}

func (p *JoinNetworkInd) Decode(bytes []byte) error {
	p.status = bytes[0]
	if p.status != LORAWAN_MSG_JOIN_NETWORK_IND_STATUS_OK && p.status != LORAWAN_MSG_JOIN_NETWORK_IND_STATUS_OK_ATTACHMENT {
		p.status = LORAWAN_MSG_JOIN_NETWORK_IND_STATUS_ERROR
	}
	if p.status == LORAWAN_MSG_JOIN_NETWORK_IND_STATUS_OK_ATTACHMENT {
		p.Address = binary.LittleEndian.Uint32(bytes[1:5])
		p.ChannelIdx = bytes[5]
		p.DataRateIdx = bytes[6]
		p.RSSI = bytes[7]
		p.SNR = bytes[8]
		p.RxSlot = bytes[9]
	}
	return nil
}

// LORAWAN_MSG_SEND_UDATA_REQ

type SendUDataReq struct {
	wimodMessageImpl
	Port    byte
	Payload []byte
}

func NewSendUDataReq(port byte, payload []byte) *SendUDataReq {
	req := &SendUDataReq{}
	req.code = LORAWAN_MSG_SEND_UDATA_REQ
	req.Port = port
	req.Payload = payload
	return req
}

func (p *SendUDataReq) String() string {
	return fmt.Sprintf("SendUDataReq[Port: %d, Payload: 0x%X]", p.Port, p.Payload)
}

func (p *SendUDataReq) Encode() ([]byte, error) {
	buff := []byte{p.Port}
	buff = append(buff, p.Payload...)
	return buff, nil
}

// LORAWAN_MSG_SEND_UDATA_RSP

type SendUDataResp struct {
	wimodMessageStatusImpl
	RemainingTime uint32
}

func NewSendUDataResp() *SendUDataResp {
	resp := &SendUDataResp{}
	resp.code = LORAWAN_MSG_SEND_UDATA_RSP
	return resp
}

func (p *SendUDataResp) String() string {
	return fmt.Sprintf("SendUDataResp[RemainingTime: %d]", p.RemainingTime)
}

func (p *SendUDataResp) Decode(payload []byte) error {
	p.status = payload[0]
	switch p.status {
	case LORAWAN_STATUS_OK:
		return nil
	case LORAWAN_STATUS_CHANNEL_BLOCKED:
		p.RemainingTime = binary.LittleEndian.Uint32(payload[1:5])
		return fmt.Errorf("LORAWAN_STATUS_CHANNEL_BLOCKED: Remaining Time: %d", p.RemainingTime)
	default:
		return lorawanStatusCheck(payload[0])
	}
}

// LORAWAN_MSG_SEND_UDATA_TX_IND

type SendUDataTxInd struct {
	wimodMessageStatusImpl
	ChannelIdx       byte
	DataRateIdx      byte
	NumTxPackets     byte
	TRXPowerLevel    byte
	RFMessageAirtime uint32
}

func NewSendUDataTxInd() *SendUDataTxInd {
	ind := &SendUDataTxInd{}
	ind.code = LORAWAN_MSG_SEND_UDATA_TX_IND
	return ind
}

func (p *SendUDataTxInd) String() string {
	return fmt.Sprintf("SendUDataTxInd[Status: 0x%02X, ChannelIdx: %d, DataRateIdx: %d, NumTxPackets: %d, TRXPowerLevel: %d, RFMessageAirtime: %d]", p.status, p.ChannelIdx, p.DataRateIdx, p.NumTxPackets, p.TRXPowerLevel, p.RFMessageAirtime)
}

func (p *SendUDataTxInd) Decode(bytes []byte) error {
	p.status = bytes[0]
	if p.status != LORAWAN_MSG_SEND_UDATA_TX_IND_STATUS_OK && p.status != LORAWAN_MSG_SEND_UDATA_TX_IND_STATUS_OK_ATTACHMENT {
		p.status = LORAWAN_MSG_SEND_UDATA_TX_IND_STATUS_ERROR
		return fmt.Errorf("LORAWAN_MSG_SEND_UDATA_TX_IND_STATUS_ERROR")
	}
	if p.status == LORAWAN_MSG_SEND_UDATA_TX_IND_STATUS_OK_ATTACHMENT {
		p.ChannelIdx = bytes[1]
		p.DataRateIdx = bytes[2]
		p.NumTxPackets = bytes[3]
		p.TRXPowerLevel = bytes[4]
		p.RFMessageAirtime = binary.LittleEndian.Uint32(bytes[5:9])
	}
	return nil
}

// LORAWAN_MSG_RECV_UDATA_IND
// LORAWAN_MSG_SEND_CDATA_REQ
// LORAWAN_MSG_SEND_CDATA_RSP
// LORAWAN_MSG_SEND_CDATA_TX_IND
// LORAWAN_MSG_RECV_CDATA_IND
// LORAWAN_MSG_RECV_ACK_IND
// LORAWAN_MSG_RECV_NO_DATA_IND
// LORAWAN_MSG_SET_RSTACK_CONFIG_REQ
// LORAWAN_MSG_SET_RSTACK_CONFIG_RSP
// LORAWAN_MSG_GET_RSTACK_CONFIG_REQ
// LORAWAN_MSG_GET_RSTACK_CONFIG_RSP
// LORAWAN_MSG_REACTIVATE_DEVICE_REQ

type ReactivateDeviceReq struct {
	wimodMessageImpl
}

func NewReactivateDeviceReq() *ReactivateDeviceReq {
	req := &ReactivateDeviceReq{}
	req.code = LORAWAN_MSG_REACTIVATE_DEVICE_REQ
	return req
}

func (p *ReactivateDeviceReq) String() string {
	return fmt.Sprintf("ReactivateDeviceReq[]")
}

func (p *ReactivateDeviceReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// LORAWAN_MSG_REACTIVATE_DEVICE_RSP

type ReactivateDeviceResp struct {
	wimodMessageStatusImpl
	Address uint32
}

func NewReactivateDeviceResp() *ReactivateDeviceResp {
	resp := &ReactivateDeviceResp{}
	resp.code = LORAWAN_MSG_REACTIVATE_DEVICE_RSP
	return resp
}

func (p *ReactivateDeviceResp) String() string {
	return fmt.Sprintf("ReactivateDeviceResp[Address: 0x%08X]", p.Address)
}

func (p *ReactivateDeviceResp) Decode(payload []byte) error {
	p.status = payload[0]
	err := lorawanStatusCheck(p.status)
	if err != nil {
		return err
	}
	p.Address = binary.LittleEndian.Uint32(payload[1:5])
	return nil
}

// LORAWAN_MSG_DEACTIVATE_DEVICE_REQ

type DeactivateDeviceReq struct {
	wimodMessageImpl
}

func NewDeactivateDeviceReq() *DeactivateDeviceReq {
	req := &DeactivateDeviceReq{}
	req.code = LORAWAN_MSG_DEACTIVATE_DEVICE_REQ
	return req
}

func (p *DeactivateDeviceReq) String() string {
	return fmt.Sprintf("DeactivateDeviceReq[]")
}

func (p *DeactivateDeviceReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// LORAWAN_MSG_DEACTIVATE_DEVICE_RSP

type DeactivateDeviceResp struct {
	wimodMessageStatusImpl
}

func NewDeactivateDeviceResp() *DeactivateDeviceResp {
	resp := &DeactivateDeviceResp{}
	resp.code = LORAWAN_MSG_DEACTIVATE_DEVICE_RSP
	return resp
}

func (p *DeactivateDeviceResp) String() string {
	return fmt.Sprintf("DeactivateDeviceResp[]")
}

func (p *DeactivateDeviceResp) Decode(payload []byte) error {
	p.status = payload[0]
	return lorawanStatusCheck(p.status)
}

// LORAWAN_MSG_FACTORY_RESET_REQ
// LORAWAN_MSG_FACTORY_RESET_RSP
// LORAWAN_MSG_SET_DEVICE_EUI_REQ
// LORAWAN_MSG_SET_DEVICE_EUI_RSP
// LORAWAN_MSG_GET_DEVICE_EUI_REQ

type GetDeviceEUIReq struct {
	wimodMessageImpl
}

func NewGetDeviceEUIReq() *GetDeviceEUIReq {
	req := &GetDeviceEUIReq{}
	req.code = LORAWAN_MSG_GET_DEVICE_EUI_REQ
	return req
}

func (p *GetDeviceEUIReq) String() string {
	return fmt.Sprintf("GetDeviceEUIReq[]")
}

func (p *GetDeviceEUIReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// LORAWAN_MSG_GET_DEVICE_EUI_RSP

type GetDeviceEUIResp struct {
	wimodMessageStatusImpl
	EUI EUI
}

func NewGetDeviceEUIResp() *GetDeviceEUIResp {
	resp := &GetDeviceEUIResp{}
	resp.code = LORAWAN_MSG_GET_DEVICE_EUI_RSP
	return resp
}

func (p *GetDeviceEUIResp) String() string {
	return fmt.Sprintf("GetDeviceEUIResp[EUI: %v]", p.EUI)
}

func (p *GetDeviceEUIResp) Decode(payload []byte) error {
	p.status = payload[0]
	err := lorawanStatusCheck(payload[0])
	if err != nil {
		return err
	}
	p.EUI = DecodeEUI(payload[1:])
	return nil
}

// LORAWAN_MSG_GET_NWK_STATUS_REQ

type GetNwkStatusReq struct {
	wimodMessageImpl
}

func NewGetNwkStatusReq() *GetNwkStatusReq {
	req := &GetNwkStatusReq{}
	req.code = LORAWAN_MSG_GET_NWK_STATUS_REQ
	return req
}

func (p *GetNwkStatusReq) String() string {
	return fmt.Sprintf("GetNwkStatusReq[]")
}

func (p *GetNwkStatusReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// LORAWAN_MSG_GET_NWK_STATUS_RSP

type GetNwkStatusResp struct {
	wimodMessageStatusImpl
	NetworkStatus  byte
	Address        uint32
	DataRateIdx    byte
	PowerLevel     byte
	MaxPayloadSize byte
}

func NewGetNwkStatusResp() *GetNwkStatusResp {
	resp := &GetNwkStatusResp{}
	resp.code = LORAWAN_MSG_GET_NWK_STATUS_RSP
	return resp
}

func (p *GetNwkStatusResp) String() string {
	return fmt.Sprintf("GetNwkStatusResp[NetworkStatus: 0x%02X, Address: 0x%08X, DataRateIdx: %d, PowerLevel: %d, MaxPayloadSize: %d]", p.NetworkStatus, p.Address, p.DataRateIdx, p.PowerLevel, p.MaxPayloadSize)
}

func (p *GetNwkStatusResp) Decode(payload []byte) error {
	p.status = payload[0]
	err := lorawanStatusCheck(p.status)
	if err != nil {
		return err
	}
	p.NetworkStatus = payload[1]
	if p.NetworkStatus == LORAWAN_NETWORK_STATUS_ACTIVE_ABP || p.NetworkStatus == LORAWAN_NETWORK_STATUS_ACTIVE_OTAA {
		p.Address = binary.LittleEndian.Uint32(payload[2:6])
		p.DataRateIdx = payload[6]
		p.PowerLevel = payload[7]
		p.MaxPayloadSize = payload[8]
	}
	return nil
}

// LORAWAN_MSG_SEND_MAC_CMD_REQ
// LORAWAN_MSG_SEND_MAC_CMD_RSP
// LORAWAN_MSG_RECV_MAC_CMD_IND
// LORAWAN_MSG_SET_CUSTOM_CFG_REQ
// LORAWAN_MSG_SET_CUSTOM_CFG_RSP
// LORAWAN_MSG_GET_CUSTOM_CFG_REQ
// LORAWAN_MSG_GET_CUSTOM_CFG_RSP
// LORAWAN_MSG_GET_SUPPORTED_BANDS_REQ
// LORAWAN_MSG_GET_SUPPORTED_BANDS_RSP
// LORAWAN_MSG_SET_LINKADRREQ_CONFIG_REQ
// LORAWAN_MSG_SET_LINKADRREQ_CONFIG_RSP
// LORAWAN_MSG_GET_LINKADRREQ_CONFIG_REQ
// LORAWAN_MSG_GET_LINKADRREQ_CONFIG_RSP
