package client

import (
	"net/rpc"
	"time"

	"github.com/enolgor/wimod-lorawan-endnode-controller/wimod"
)

type WimodClient struct {
	Client *rpc.Client
}

// Ping

func (c *WimodClient) Ping() error {
	resp := 0
	return c.Client.Call("WimodServer.Ping", 0, &resp)
}

// GetDeviceInfo

func (c *WimodClient) GetDeviceInfo() (*wimod.GetDeviceInfoResp, error) {
	resp := wimod.NewGetDeviceInfoResp()
	err := c.Client.Call("WimodServer.GetDeviceInfo", 0, resp)
	return resp, err
}

// GetFWInfo

func (c *WimodClient) GetFWInfo() (*wimod.GetFWInfoResp, error) {
	resp := wimod.NewGetFWInfoResp()
	err := c.Client.Call("WimodServer.GetFWInfo", 0, resp)
	return resp, err
}

//  Reset

func (c *WimodClient) Reset() error {
	resp := 0
	return c.Client.Call("WimodServer.Reset", 0, &resp)
}

// SetOPMode

func (c *WimodClient) SetOPMode(mode byte) error {
	resp := 0
	return c.Client.Call("WimodServer.SetOPMode", wimod.NewSetOPModeReq(mode), &resp)
}

// GetOPMode

func (c *WimodClient) GetOPMode() (*wimod.GetOPModeResp, error) {
	resp := wimod.NewGetOPModeResp()
	err := c.Client.Call("WimodServer.GetOPMode", 0, resp)
	return resp, err
}

// SetRTC

func (c *WimodClient) SetRTC(time time.Time) error {
	resp := 0
	return c.Client.Call("WimodServer.SetRTC", wimod.NewSetRTCReq(time), &resp)
}

// GetRTC

func (c *WimodClient) GetRTC() (*wimod.GetRTCResp, error) {
	resp := wimod.NewGetRTCResp()
	err := c.Client.Call("WimodServer.GetRTC", 0, resp)
	return resp, err
}

// GetDeviceStatus

func (c *WimodClient) GetDeviceStatus() (*wimod.GetDeviceStatusResp, error) {
	resp := wimod.NewGetDeviceStatusResp()
	err := c.Client.Call("WimodServer.GetDeviceStatus", 0, resp)
	return resp, err
}

// SetRTCAlarm

func (c *WimodClient) SetRTCAlarm(alarmType, hour, minutes, seconds byte) error {
	resp := 0
	return c.Client.Call("WimodServer.SetRTCAlarm", wimod.NewSetRTCAlarmReq(alarmType, hour, minutes, seconds), &resp)
}

// ClearRTCAlarm

func (c *WimodClient) ClearRTCAlarm() error {
	resp := 0
	return c.Client.Call("WimodServer.ClearRTCAlarm", 0, &resp)
}

// GetRTCAlarm

func (c *WimodClient) GetRTCAlarm() (*wimod.GetRTCAlarmResp, error) {
	resp := wimod.NewGetRTCAlarmResp()
	err := c.Client.Call("WimodServer.GetRTCAlarm", 0, resp)
	return resp, err
}

// RTCAlarmInd

func (c *WimodClient) RTCAlarmInd() (*wimod.RTCAlarmInd, error) {
	ind := wimod.NewRTCAlarmInd()
	err := c.Client.Call("WimodServer.RTCAlarmInd", 0, ind)
	return ind, err
}

// ActivateDevice

func (c *WimodClient) ActivateDevice(address uint32, appSessKey wimod.Key, nwkSessKey wimod.Key) error {
	resp := 0
	return c.Client.Call("WimodServer.ActivateDevice", wimod.NewActivateDeviceReq(address, appSessKey, nwkSessKey), &resp)
}

// SetJoinParam

func (c *WimodClient) SetJoinParam(appEUI wimod.EUI, appKey wimod.Key) error {
	resp := 0
	return c.Client.Call("WimodServer.SetJoinParam", wimod.NewSetJoinParamReq(appEUI, appKey), &resp)
}

// JoinNetwork

func (c *WimodClient) JoinNetwork() error {
	resp := 0
	return c.Client.Call("WimodServer.JoinNetwork", 0, &resp)
}

// JoinNetworkTxInd

func (c *WimodClient) JoinNetworkTxInd() (*wimod.JoinNetworkTxInd, error) {
	ind := wimod.NewJoinNetworkTxInd()
	err := c.Client.Call("WimodServer.JoinNetworkTxInd", 0, ind)
	return ind, err
}

// JoinNetworkInd

func (c *WimodClient) JoinNetworkInd() (*wimod.JoinNetworkInd, error) {
	ind := wimod.NewJoinNetworkInd()
	err := c.Client.Call("WimodServer.JoinNetworkInd", 0, ind)
	return ind, err
}

// SendUData

func (c *WimodClient) SendUData(port byte, payload []byte) (*wimod.SendUDataResp, error) {
	resp := wimod.NewSendUDataResp()
	err := c.Client.Call("WimodServer.SendUData", wimod.NewSendUDataReq(port, payload), resp)
	return resp, err
}

// SendUDataTxInd

func (c *WimodClient) SendUDataTxInd() (*wimod.SendUDataTxInd, error) {
	ind := wimod.NewSendUDataTxInd()
	err := c.Client.Call("WimodServer.SendUDataTxInd", 0, ind)
	return ind, err
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

// GetRStackConfig

func (c *WimodClient) GetRStackConfig() (*wimod.GetRStackConfigResp, error) {
	resp := wimod.NewGetRStackConfigResp()
	err := c.Client.Call("WimodServer.GetRStackConfig", 0, resp)
	return resp, err
}

// ReactivateDevice

func (c *WimodClient) ReactivateDevice() (*wimod.ReactivateDeviceResp, error) {
	resp := wimod.NewReactivateDeviceResp()
	err := c.Client.Call("WimodServer.ReactivateDevice", 0, resp)
	return resp, err
}

// DeactivateDevice

func (c *WimodClient) DeactivateDevice() error {
	resp := 0
	return c.Client.Call("WimodServer.DeactivateDevice", 0, &resp)
}

// LORAWAN_MSG_FACTORY_RESET_REQ
// LORAWAN_MSG_FACTORY_RESET_RSP
// LORAWAN_MSG_SET_DEVICE_EUI_REQ
// LORAWAN_MSG_SET_DEVICE_EUI_RSP
// LORAWAN_MSG_GET_DEVICE_EUI_REQ

// GetDeviceEUI

func (c *WimodClient) GetDeviceEUI() (*wimod.GetDeviceEUIResp, error) {
	resp := wimod.NewGetDeviceEUIResp()
	err := c.Client.Call("WimodServer.GetDeviceEUI", 0, resp)
	return resp, err
}

// GetNwkStatus

func (c *WimodClient) GetNwkStatus() (*wimod.GetNwkStatusResp, error) {
	resp := wimod.NewGetNwkStatusResp()
	err := c.Client.Call("WimodServer.GetNwkStatus", 0, resp)
	return resp, err
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
