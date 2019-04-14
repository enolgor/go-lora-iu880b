package wimod

import (
	"encoding/binary"
	"fmt"
	"time"
)

// DEVMGMT_MSG_PING_REQ

type PingReq struct {
	wimodmessage
}

func NewPingReq() *PingReq {
	req := &PingReq{}
	req.dst = DEVMGMT_ID
	req.id = DEVMGMT_MSG_PING_REQ
	return req
}

func (p *PingReq) String() string {
	return fmt.Sprintf("PingReq[]")
}

func (p *PingReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_PING_RSP

type PingResp struct {
	wimodmessage
}

func NewPingResp() *PingResp {
	resp := &PingResp{}
	resp.dst = DEVMGMT_ID
	resp.id = DEVMGMT_MSG_PING_RSP
	return resp
}

func (p *PingResp) String() string {
	return fmt.Sprintf("PingResp[]")
}

func (p *PingResp) Decode(bytes []byte) error {
	return nil
}

// DEVMGMT_MSG_GET_DEVICE_INFO_REQ

type GetDeviceInfoReq struct {
	wimodmessage
}

func NewGetDeviceInfoReq() *GetDeviceInfoReq {
	req := &GetDeviceInfoReq{}
	req.dst = DEVMGMT_ID
	req.id = DEVMGMT_MSG_GET_DEVICE_INFO_REQ
	return req
}

func (p *GetDeviceInfoReq) String() string {
	return fmt.Sprintf("GetDeviceInfoReq[]")
}

func (p *GetDeviceInfoReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_DEVICE_INFO_RSP

type GetDeviceInfoResp struct {
	wimodmessage
	ModuleType    byte
	DeviceAddress uint32
	DeviceID      uint32
}

func NewGetDeviceInfoResp() *GetDeviceInfoResp {
	resp := &GetDeviceInfoResp{}
	resp.dst = DEVMGMT_ID
	resp.id = DEVMGMT_MSG_GET_DEVICE_INFO_RSP
	return resp
}

func (p *GetDeviceInfoResp) String() string {
	return fmt.Sprintf("GetDeviceInfoResp[ModuleType: 0x%X, DeviceAddress: 0x%X, DeviceID: 0x%X]", p.ModuleType, p.DeviceAddress, p.DeviceID)
}

func (p *GetDeviceInfoResp) Decode(bytes []byte) error {
	p.ModuleType = bytes[0]
	p.DeviceAddress = binary.LittleEndian.Uint32(bytes[1:5])
	p.DeviceID = binary.LittleEndian.Uint32(bytes[5:9])
	return nil
}

// DEVMGMT_MSG_GET_FW_INFO_REQ

type GetFWInfoReq struct {
	wimodmessage
}

func NewGetFWInfoReq() *GetFWInfoReq {
	req := &GetFWInfoReq{}
	req.dst = DEVMGMT_ID
	req.id = DEVMGMT_MSG_GET_FW_INFO_REQ
	return req
}

func (p *GetFWInfoReq) String() string {
	return fmt.Sprintf("GetFWInfoReq[]")
}

func (p *GetFWInfoReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_FW_INFO_RSP

type GetFWInfoResp struct {
	wimodmessage
	MinorVersion  byte
	MajorVersion  byte
	Build         uint16
	BuildDate     string
	FirmwareImage string
}

func NewGetFWInfoResp() *GetFWInfoResp {
	resp := &GetFWInfoResp{}
	resp.dst = DEVMGMT_ID
	resp.id = DEVMGMT_MSG_GET_FW_INFO_RSP
	return resp
}

func (p *GetFWInfoResp) String() string {
	return fmt.Sprintf("GetFWInfoResp[MinorVersion: %d, MajorVersion: %d, Build: %d, BuildDate: %s, FirmwareImage: %s]", p.MinorVersion, p.MajorVersion, p.Build, p.BuildDate, p.FirmwareImage)
}

func (p *GetFWInfoResp) Decode(bytes []byte) error {
	p.MinorVersion = bytes[0]
	p.MajorVersion = bytes[1]
	p.Build = binary.LittleEndian.Uint16(bytes[2:4])
	p.BuildDate = string(bytes[4:14])
	p.FirmwareImage = string(bytes[14:])
	return nil
}

// DEVMGMT_MSG_RESET_REQ

// DEVMGMT_MSG_RESET_RSP

// DEVMGMT_MSG_SET_OPMODE_REQ

// DEVMGMT_MSG_SET_OPMODE_RSP

// DEVMGMT_MSG_GET_OPMODE_REQ

// DEVMGMT_MSG_GET_OPMODE_RSP

// DEVMGMT_MSG_SET_RTC_REQ

type SetRTCReq struct {
	wimodmessage
	Time time.Time
}

func NewSetRTCReq() *SetRTCReq {
	req := &SetRTCReq{}
	req.dst = DEVMGMT_ID
	req.id = DEVMGMT_MSG_SET_RTC_REQ
	return req
}

func (p *SetRTCReq) String() string {
	return fmt.Sprintf("SetRTCReq[Time: %s]", p.Time)
}

func (p *SetRTCReq) Encode() ([]byte, error) {
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, EncodeRTCTime(p.Time))
	return buff, nil
}

// DEVMGMT_MSG_SET_RTC_RSP

type SetRTCResp struct {
	wimodmessage
}

func NewSetRTCResp() *SetRTCResp {
	resp := &SetRTCResp{}
	resp.dst = DEVMGMT_ID
	resp.id = DEVMGMT_MSG_SET_RTC_RSP
	return resp
}

func (p *SetRTCResp) String() string {
	return fmt.Sprintf("SetRTCResp[]")
}

func (p *SetRTCResp) Decode(bytes []byte) error {
	return nil
}

// DEVMGMT_MSG_GET_RTC_REQ

type GetRTCReq struct {
	wimodmessage
}

func NewGetRTCReq() *GetRTCReq {
	req := &GetRTCReq{}
	req.dst = DEVMGMT_ID
	req.id = DEVMGMT_MSG_GET_RTC_REQ
	return req
}

func (p *GetRTCReq) String() string {
	return fmt.Sprintf("GetRTCReq[]")
}

func (p *GetRTCReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_RTC_RSP

type GetRTCResp struct {
	wimodmessage
	Time time.Time
}

func NewGetRTCResp() *GetRTCResp {
	resp := &GetRTCResp{}
	resp.dst = DEVMGMT_ID
	resp.id = DEVMGMT_MSG_GET_RTC_RSP
	return resp
}

func (p *GetRTCResp) String() string {
	return fmt.Sprintf("GetRTCResp[Time: %s]", p.Time)
}

func (p *GetRTCResp) Decode(bytes []byte) error {
	p.Time = DecodeRTCTime(binary.LittleEndian.Uint32(bytes))
	return nil
}

// DEVMGMT_MSG_GET_DEVICE_STATUS_REQ

type GetDeviceStatusReq struct {
	wimodmessage
}

func NewGetDeviceStatusReq() *GetDeviceStatusReq {
	req := &GetDeviceStatusReq{}
	req.dst = DEVMGMT_ID
	req.id = DEVMGMT_MSG_GET_DEVICE_STATUS_REQ
	return req
}

func (p *GetDeviceStatusReq) String() string {
	return fmt.Sprintf("GetDeviceStatusReq[]")
}

func (p *GetDeviceStatusReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_DEVICE_STATUS_RSP

type GetDeviceStatusResp struct {
	wimodmessage
	SystemTickResolution byte
	SystemTicks          uint32
	TargetTime           time.Time
	NVMStatus            uint16
	BatteryLevel         uint16
	ExtraStatus          uint16
	TxUData              uint32
	TxCData              uint32
	TxError              uint32
	Rx1UData             uint32
	Rx1CData             uint32
	Rx1MICError          uint32
	Rx2UData             uint32
	Rx2CData             uint32
	Rx2MICError          uint32
	TxJoin               uint32
	RxAccept             uint32
}

func NewGetDeviceStatusResp() *GetDeviceStatusResp {
	resp := &GetDeviceStatusResp{}
	resp.dst = DEVMGMT_ID
	resp.id = DEVMGMT_MSG_GET_DEVICE_STATUS_RSP
	return resp
}

func (p *GetDeviceStatusResp) String() string {
	return fmt.Sprintf("GetDeviceStatusResp[SystemTickResolution: %dms, SystemTicks: %d, TargetTime: %s, NVMStatus: %016bb, BatteryLevel: %dmV, ExtraStatus: 0x%04X, TxU-Data: %d, TxC-Data: %d, TxError: %d, Rx1U-Data: %d, Rx1C-Data: %d, Rx1MIC-Error: %d, Rx2U-Data: %d, Rx2C-Data: %d, Rx2MIC-Error: %d, TxJoin: %d, RxAccept: %d]", p.SystemTickResolution, p.SystemTicks, p.TargetTime, p.NVMStatus, p.BatteryLevel, p.ExtraStatus, p.TxUData, p.TxCData, p.TxError, p.Rx1UData, p.Rx1CData, p.Rx1MICError, p.Rx2UData, p.Rx2CData, p.Rx2MICError, p.TxJoin, p.RxAccept)
}

func (p *GetDeviceStatusResp) Decode(bytes []byte) error {
	p.SystemTickResolution = bytes[0]
	p.SystemTicks = binary.LittleEndian.Uint32(bytes[1:5])
	p.TargetTime = DecodeRTCTime(binary.LittleEndian.Uint32(bytes[5:9]))
	p.NVMStatus = binary.LittleEndian.Uint16(bytes[9:11])
	p.BatteryLevel = binary.LittleEndian.Uint16(bytes[11:13])
	p.ExtraStatus = binary.LittleEndian.Uint16(bytes[13:15])
	p.TxUData = binary.LittleEndian.Uint32(bytes[15:19])
	p.TxCData = binary.LittleEndian.Uint32(bytes[19:23])
	p.TxError = binary.LittleEndian.Uint32(bytes[23:27])
	p.Rx1UData = binary.LittleEndian.Uint32(bytes[27:31])
	p.Rx1CData = binary.LittleEndian.Uint32(bytes[31:35])
	p.Rx1MICError = binary.LittleEndian.Uint32(bytes[35:39])
	p.Rx2UData = binary.LittleEndian.Uint32(bytes[39:43])
	p.Rx2CData = binary.LittleEndian.Uint32(bytes[43:47])
	p.Rx2MICError = binary.LittleEndian.Uint32(bytes[47:51])
	p.TxJoin = binary.LittleEndian.Uint32(bytes[51:55])
	p.RxAccept = binary.LittleEndian.Uint32(bytes[55:59])
	return nil
}

// DEVMGMT_MSG_SET_RTC_ALARM_REQ
// DEVMGMT_MSG_SET_RTC_ALARM_RSP
// DEVMGMT_MSG_CLEAR_RTC_ALARM_REQ
// DEVMGMT_MSG_CLEAR_RTC_ALARM_RSP
// DEVMGMT_MSG_GET_RTC_ALARM_REQ
// DEVMGMT_MSG_GET_RTC_ALARM_RSP
// DEVMGMT_MSG_RTC_ALARM_IND
