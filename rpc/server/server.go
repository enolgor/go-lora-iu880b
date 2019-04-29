package server

import (
	"github.com/enolgor/wimod-lorawan-endnode-controller/controller"
	"github.com/enolgor/wimod-lorawan-endnode-controller/wimod"
)

type WimodServer struct {
	Controller *controller.WiModController
}

// Ping

func (s *WimodServer) Ping(_ *int, _ *int) error {
	return s.Controller.Request(wimod.NewPingReq(), wimod.NewPingResp())
}

// GetDeviceInfo

func (s *WimodServer) GetDeviceInfo(_ *int, response *wimod.GetDeviceInfoResp) error {
	return s.Controller.Request(wimod.NewGetDeviceInfoReq(), response)
}

// GetFWInfo

func (s *WimodServer) GetFWInfo(_ *int, response *wimod.GetFWInfoResp) error {
	return s.Controller.Request(wimod.NewGetFWInfoReq(), response)
}

//  Reset

func (s *WimodServer) Reset(_ *int, _ *int) error {
	return s.Controller.Request(wimod.NewResetReq(), wimod.NewResetResp())
}

// SetOPMode

func (s *WimodServer) SetOPMode(request *wimod.SetOPModeReq, _ *int) error {
	return s.Controller.Request(request, wimod.NewSetOPModeResp())
}

// GetOPMode

func (s *WimodServer) GetOPMode(_ *int, response *wimod.GetOPModeResp) error {
	return s.Controller.Request(wimod.NewGetOPModeReq(), response)
}

// SetRTC

func (s *WimodServer) SetRTC(request *wimod.SetRTCReq, _ *int) error {
	return s.Controller.Request(request, wimod.NewSetRTCResp())
}

// GetRTC

func (s *WimodServer) GetRTC(_ *int, response *wimod.GetRTCResp) error {
	return s.Controller.Request(wimod.NewGetRTCReq(), response)
}

// GetDeviceStatus

func (s *WimodServer) GetDeviceStatus(_ *int, response *wimod.GetDeviceStatusResp) error {
	return s.Controller.Request(wimod.NewGetDeviceStatusReq(), response)
}

// SetRTCAlarm

func (s *WimodServer) SetRTCAlarm(request *wimod.SetRTCAlarmReq, _ *int) error {
	return s.Controller.Request(request, wimod.NewSetRTCAlarmResp())
}

// ClearRTCAlarm

func (s *WimodServer) ClearRTCAlarm(_ *int, _ *int) error {
	return s.Controller.Request(wimod.NewClearRTCAlarmReq(), wimod.NewClearRTCAlarmResp())
}

// GetRTCAlarm

func (s *WimodServer) GetRTCAlarm(_ *int, response *wimod.GetRTCAlarmResp) error {
	return s.Controller.Request(wimod.NewGetRTCAlarmReq(), response)
}

// RTCAlarmInd

func (s *WimodServer) RTCAlarmInd(_ *int, ind *wimod.RTCAlarmInd) error {
	return s.Controller.ReadSpecificInd(ind)
}

// ActivateDevice

func (s *WimodServer) ActivateDevice(request *wimod.ActivateDeviceReq, _ *int) error {
	return s.Controller.Request(request, wimod.NewActivateDeviceResp())
}

// SetJoinParam

func (s *WimodServer) SetJoinParam(request *wimod.SetJoinParamReq, _ *int) error {
	return s.Controller.Request(request, wimod.NewSetJoinParamResp())
}

// JoinNetwork

func (s *WimodServer) JoinNetwork(_ *int, _ *int) error {
	return s.Controller.Request(wimod.NewJoinNetworkReq(), wimod.NewJoinNetworkResp())
}

// JoinNetworkTxInd

func (s *WimodServer) JoinNetworkTxInd(_ *int, ind *wimod.JoinNetworkTxInd) error {
	return s.Controller.ReadSpecificInd(ind)
}

// JoinNetworkInd

func (s *WimodServer) JoinNetworkInd(_ *int, ind *wimod.JoinNetworkInd) error {
	return s.Controller.ReadSpecificInd(ind)
}

// SendUData

func (s *WimodServer) SendUData(request *wimod.SendUDataReq, response *wimod.SendUDataResp) error {
	return s.Controller.Request(request, response)
}

// SendUDataTxInd

func (s *WimodServer) SendUDataTxInd(_ *int, ind *wimod.SendUDataTxInd) error {
	return s.Controller.ReadSpecificInd(ind)
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

func (s *WimodServer) GetRStackConfig(_ *int, response *wimod.GetRStackConfigResp) error {
	return s.Controller.Request(wimod.NewGetRStackConfigReq(), response)
}

// ReactivateDevice

func (s *WimodServer) ReactivateDevice(_ *int, response *wimod.ReactivateDeviceResp) error {
	return s.Controller.Request(wimod.NewReactivateDeviceReq(), response)
}

// DeactivateDevice

func (s *WimodServer) DeactivateDevice(_ *int, _ *int) error {
	return s.Controller.Request(wimod.NewDeactivateDeviceReq(), wimod.NewDeactivateDeviceResp())
}

// LORAWAN_MSG_FACTORY_RESET_REQ
// LORAWAN_MSG_FACTORY_RESET_RSP
// LORAWAN_MSG_SET_DEVICE_EUI_REQ
// LORAWAN_MSG_SET_DEVICE_EUI_RSP
// LORAWAN_MSG_GET_DEVICE_EUI_REQ

// GetDeviceEUI

func (s *WimodServer) GetDeviceEUI(_ *int, response *wimod.GetDeviceEUIResp) error {
	return s.Controller.Request(wimod.NewGetDeviceEUIReq(), response)
}

// GetNwkStatus

func (s *WimodServer) GetNwkStatus(_ *int, response *wimod.GetNwkStatusResp) error {
	return s.Controller.Request(wimod.NewGetNwkStatusReq(), response)
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
