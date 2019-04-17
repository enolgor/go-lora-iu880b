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
	return []byte{}, nil
}

// LORAWAN_MSG_SET_JOIN_PARAM_RSP

type SetJoinParamResp struct {
	wimodMessageImpl
}

func NewSetJoinParamResp() *SetJoinParamResp {
	resp := &SetJoinParamResp{}
	resp.code = LORAWAN_MSG_SET_JOIN_PARAM_RSP
	return resp
}

func (p *SetJoinParamResp) String() string {
	return fmt.Sprintf("SetJoinParamResp[]")
}

func (p *SetJoinParamResp) Decode(bytes []byte) error {
	return nil
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
	wimodMessageImpl
}

func NewJoinNetworkResp() *JoinNetworkResp {
	resp := &JoinNetworkResp{}
	resp.code = LORAWAN_MSG_JOIN_NETWORK_RSP
	return resp
}

func (p *JoinNetworkResp) String() string {
	return fmt.Sprintf("JoinNetworkResp[]")
}

func (p *JoinNetworkResp) Decode(bytes []byte) error {
	return nil
}

// LORAWAN_MSG_JOIN_NETWORK_TX_IND

type JoinNetworkTxInd struct {
	wimodMessageIndImpl
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
	return fmt.Sprintf("JoinNetworkTxInd[ChannelIdx: %d, DataRateIdx: %d, NumTxPackets: %d, TRXPowerLevel: %d, RFMessageAirtime: %d]", p.ChannelIdx, p.DataRateIdx, p.NumTxPackets, p.TRXPowerLevel, p.RFMessageAirtime)
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
// LORAWAN_MSG_SEND_UDATA_REQ
// LORAWAN_MSG_SEND_UDATA_RSP
// LORAWAN_MSG_SEND_UDATA_TX_IND
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
	wimodMessageImpl
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

func (p *ReactivateDeviceResp) Decode(bytes []byte) error {
	p.Address = binary.LittleEndian.Uint32(bytes[:4])
	return nil
}

// LORAWAN_MSG_DEACTIVATE_DEVICE_REQ
// LORAWAN_MSG_DEACTIVATE_DEVICE_RSP
// LORAWAN_MSG_FACTORY_RESET_REQ
// LORAWAN_MSG_FACTORY_RESET_RSP
// LORAWAN_MSG_SET_DEVICE_EUI_REQ
// LORAWAN_MSG_SET_DEVICE_EUI_RSP
// LORAWAN_MSG_GET_DEVICE_EUI_REQ
// LORAWAN_MSG_GET_DEVICE_EUI_RSP
// LORAWAN_MSG_GET_NWK_STATUS_REQ
// LORAWAN_MSG_GET_NWK_STATUS_RSP
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
