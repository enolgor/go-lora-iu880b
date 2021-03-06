package wimod

import (
	"encoding/binary"
	"fmt"
	"time"
)

// DEVMGMT_MSG_PING_REQ

type PingReq struct {
	wimodMessageImpl
}

func NewPingReq() *PingReq {
	req := &PingReq{}
	req.Init()
	return req
}

func (p *PingReq) Init() {
	p.code = DEVMGMT_MSG_PING_REQ
}

func (p *PingReq) String() string {
	return fmt.Sprintf("PingReq[]")
}

func (p *PingReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_PING_RSP

type PingResp struct {
	wimodMessageStatusImpl
}

func NewPingResp() *PingResp {
	resp := &PingResp{}
	resp.Init()
	return resp
}

func (p *PingResp) Init() {
	p.code = DEVMGMT_MSG_PING_RSP
}

func (p *PingResp) String() string {
	return fmt.Sprintf("PingResp[]")
}

func (p *PingResp) Decode(payload []byte) error {
	p.Status = payload[0]
	return devMgmtStatusCheck(p.Status)
}

// DEVMGMT_MSG_GET_DEVICE_INFO_REQ

type GetDeviceInfoReq struct {
	wimodMessageImpl
}

func NewGetDeviceInfoReq() *GetDeviceInfoReq {
	req := &GetDeviceInfoReq{}
	req.Init()
	return req
}

func (p *GetDeviceInfoReq) Init() {
	p.code = DEVMGMT_MSG_GET_DEVICE_INFO_REQ
}

func (p *GetDeviceInfoReq) String() string {
	return fmt.Sprintf("GetDeviceInfoReq[]")
}

func (p *GetDeviceInfoReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_DEVICE_INFO_RSP

type GetDeviceInfoResp struct {
	wimodMessageStatusImpl
	ModuleType    byte
	DeviceAddress uint32
	DeviceID      uint32
}

func NewGetDeviceInfoResp() *GetDeviceInfoResp {
	resp := &GetDeviceInfoResp{}
	resp.Init()
	return resp
}

func (p *GetDeviceInfoResp) Init() {
	p.code = DEVMGMT_MSG_GET_DEVICE_INFO_RSP
}

func (p *GetDeviceInfoResp) String() string {
	return fmt.Sprintf("GetDeviceInfoResp[ModuleType: 0x%X, DeviceAddress: 0x%X, DeviceID: 0x%X]", p.ModuleType, p.DeviceAddress, p.DeviceID)
}

func (p *GetDeviceInfoResp) Decode(payload []byte) error {
	p.Status = payload[0]
	err := devMgmtStatusCheck(p.Status)
	if err != nil {
		return err
	}
	bytes := payload[1:]
	p.ModuleType = bytes[0]
	p.DeviceAddress = binary.LittleEndian.Uint32(bytes[1:5])
	p.DeviceID = binary.LittleEndian.Uint32(bytes[5:9])
	return nil
}

// DEVMGMT_MSG_GET_FW_INFO_REQ

type GetFWInfoReq struct {
	wimodMessageImpl
}

func NewGetFWInfoReq() *GetFWInfoReq {
	req := &GetFWInfoReq{}
	req.Init()
	return req
}

func (p *GetFWInfoReq) Init() {
	p.code = DEVMGMT_MSG_GET_FW_INFO_REQ
}

func (p *GetFWInfoReq) String() string {
	return fmt.Sprintf("GetFWInfoReq[]")
}

func (p *GetFWInfoReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_FW_INFO_RSP

type GetFWInfoResp struct {
	wimodMessageStatusImpl
	MinorVersion  byte
	MajorVersion  byte
	Build         uint16
	BuildDate     string
	FirmwareImage string
}

func NewGetFWInfoResp() *GetFWInfoResp {
	resp := &GetFWInfoResp{}
	resp.Init()
	return resp
}

func (p *GetFWInfoResp) Init() {
	p.code = DEVMGMT_MSG_GET_FW_INFO_RSP
}

func (p *GetFWInfoResp) String() string {
	return fmt.Sprintf("GetFWInfoResp[MinorVersion: %d, MajorVersion: %d, Build: %d, BuildDate: %s, FirmwareImage: %s]", p.MinorVersion, p.MajorVersion, p.Build, p.BuildDate, p.FirmwareImage)
}

func (p *GetFWInfoResp) Decode(payload []byte) error {
	p.Status = payload[0]
	err := devMgmtStatusCheck(p.Status)
	if err != nil {
		return err
	}
	bytes := payload[1:]
	p.MinorVersion = bytes[0]
	p.MajorVersion = bytes[1]
	p.Build = binary.LittleEndian.Uint16(bytes[2:4])
	p.BuildDate = string(bytes[4:14])
	p.FirmwareImage = string(bytes[14:])
	return nil
}

// DEVMGMT_MSG_RESET_REQ

type ResetReq struct {
	wimodMessageImpl
}

func NewResetReq() *ResetReq {
	req := &ResetReq{}
	req.Init()
	return req
}

func (p *ResetReq) Init() {
	p.code = DEVMGMT_MSG_RESET_REQ
}

func (p *ResetReq) String() string {
	return fmt.Sprintf("ResetReq[]")
}

func (p *ResetReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_RESET_RSP

type ResetResp struct {
	wimodMessageStatusImpl
}

func NewResetResp() *ResetResp {
	resp := &ResetResp{}
	resp.Init()
	return resp
}

func (p *ResetResp) Init() {
	p.code = DEVMGMT_MSG_RESET_RSP
}

func (p *ResetResp) String() string {
	return fmt.Sprintf("ResetResp[]")
}

func (p *ResetResp) Decode(payload []byte) error {
	p.Status = payload[0]
	return devMgmtStatusCheck(p.Status)
}

// DEVMGMT_MSG_SET_OPMODE_REQ

type SetOPModeReq struct {
	wimodMessageImpl
	Mode byte
}

func NewSetOPModeReq(mode byte) *SetOPModeReq {
	req := &SetOPModeReq{}
	req.Init()
	req.Mode = mode
	return req
}

func (p *SetOPModeReq) Init() {
	p.code = DEVMGMT_MSG_SET_OPMODE_REQ
}

func (p *SetOPModeReq) String() string {
	return fmt.Sprintf("SetOPModeReq[Mode: %02X]", p.Mode)
}

func (p *SetOPModeReq) Encode() ([]byte, error) {
	return []byte{p.Mode}, nil
}

// DEVMGMT_MSG_SET_OPMODE_RSP

type SetOPModeResp struct {
	wimodMessageStatusImpl
}

func NewSetOPModeResp() *SetOPModeResp {
	resp := &SetOPModeResp{}
	resp.Init()
	return resp
}

func (p *SetOPModeResp) Init() {
	p.code = DEVMGMT_MSG_SET_OPMODE_RSP
}

func (p *SetOPModeResp) String() string {
	return fmt.Sprintf("SetOPModeResp[]")
}

func (p *SetOPModeResp) Decode(payload []byte) error {
	p.Status = payload[0]
	return devMgmtStatusCheck(p.Status)
}

// DEVMGMT_MSG_GET_OPMODE_REQ

type GetOPModeReq struct {
	wimodMessageImpl
}

func NewGetOPModeReq() *GetOPModeReq {
	req := &GetOPModeReq{}
	req.Init()
	return req
}

func (p *GetOPModeReq) Init() {
	p.code = DEVMGMT_MSG_GET_OPMODE_REQ
}

func (p *GetOPModeReq) String() string {
	return fmt.Sprintf("GetOPModeReq[]")
}

func (p *GetOPModeReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_OPMODE_RSP

type GetOPModeResp struct {
	wimodMessageStatusImpl
	Mode byte
}

func NewGetOPModeResp() *GetOPModeResp {
	resp := &GetOPModeResp{}
	resp.Init()
	return resp
}

func (p *GetOPModeResp) Init() {
	p.code = DEVMGMT_MSG_GET_OPMODE_RSP
}

func (p *GetOPModeResp) String() string {
	return fmt.Sprintf("GetOPModeResp[Mode: %02X]", p.Mode)
}

func (p *GetOPModeResp) Decode(payload []byte) error {
	p.Status = payload[0]
	err := devMgmtStatusCheck(p.Status)
	if err != nil {
		return err
	}
	p.Mode = payload[1]
	return nil
}

// DEVMGMT_MSG_SET_RTC_REQ

type SetRTCReq struct {
	wimodMessageImpl
	Time time.Time
}

func NewSetRTCReq(time time.Time) *SetRTCReq {
	req := &SetRTCReq{}
	req.Init()
	req.Time = time
	return req
}

func (p *SetRTCReq) Init() {
	p.code = DEVMGMT_MSG_SET_RTC_REQ
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
	wimodMessageStatusImpl
}

func NewSetRTCResp() *SetRTCResp {
	resp := &SetRTCResp{}
	resp.Init()
	return resp
}

func (p *SetRTCResp) Init() {
	p.code = DEVMGMT_MSG_SET_RTC_RSP
}

func (p *SetRTCResp) String() string {
	return fmt.Sprintf("SetRTCResp[]")
}

func (p *SetRTCResp) Decode(payload []byte) error {
	p.Status = payload[0]
	return devMgmtStatusCheck(p.Status)
}

// DEVMGMT_MSG_GET_RTC_REQ

type GetRTCReq struct {
	wimodMessageImpl
}

func NewGetRTCReq() *GetRTCReq {
	req := &GetRTCReq{}
	req.Init()
	return req
}

func (p *GetRTCReq) Init() {
	p.code = DEVMGMT_MSG_GET_RTC_REQ
}

func (p *GetRTCReq) String() string {
	return fmt.Sprintf("GetRTCReq[]")
}

func (p *GetRTCReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_RTC_RSP

type GetRTCResp struct {
	wimodMessageStatusImpl
	Time time.Time
}

func NewGetRTCResp() *GetRTCResp {
	resp := &GetRTCResp{}
	resp.Init()
	return resp
}

func (p *GetRTCResp) Init() {
	p.code = DEVMGMT_MSG_GET_RTC_RSP
}

func (p *GetRTCResp) String() string {
	return fmt.Sprintf("GetRTCResp[Time: %s]", p.Time)
}

func (p *GetRTCResp) Decode(payload []byte) error {
	p.Status = payload[0]
	err := devMgmtStatusCheck(p.Status)
	if err != nil {
		return err
	}
	p.Time = DecodeRTCTime(binary.LittleEndian.Uint32(payload[1:]))
	return nil
}

// DEVMGMT_MSG_GET_DEVICE_STATUS_REQ

type GetDeviceStatusReq struct {
	wimodMessageImpl
}

func NewGetDeviceStatusReq() *GetDeviceStatusReq {
	req := &GetDeviceStatusReq{}
	req.Init()
	return req
}

func (p *GetDeviceStatusReq) Init() {
	p.code = DEVMGMT_MSG_GET_DEVICE_STATUS_REQ
}

func (p *GetDeviceStatusReq) String() string {
	return fmt.Sprintf("GetDeviceStatusReq[]")
}

func (p *GetDeviceStatusReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_DEVICE_STATUS_RSP

type GetDeviceStatusResp struct {
	wimodMessageStatusImpl
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
	resp.Init()
	return resp
}

func (p *GetDeviceStatusResp) Init() {
	p.code = DEVMGMT_MSG_GET_DEVICE_STATUS_RSP
}

func (p *GetDeviceStatusResp) String() string {
	return fmt.Sprintf("GetDeviceStatusResp[SystemTickResolution: %dms, SystemTicks: %d, TargetTime: %s, NVMStatus: %016bb, BatteryLevel: %dmV, ExtraStatus: 0x%04X, TxU-Data: %d, TxC-Data: %d, TxError: %d, Rx1U-Data: %d, Rx1C-Data: %d, Rx1MIC-Error: %d, Rx2U-Data: %d, Rx2C-Data: %d, Rx2MIC-Error: %d, TxJoin: %d, RxAccept: %d]", p.SystemTickResolution, p.SystemTicks, p.TargetTime, p.NVMStatus, p.BatteryLevel, p.ExtraStatus, p.TxUData, p.TxCData, p.TxError, p.Rx1UData, p.Rx1CData, p.Rx1MICError, p.Rx2UData, p.Rx2CData, p.Rx2MICError, p.TxJoin, p.RxAccept)
}

func (p *GetDeviceStatusResp) Decode(payload []byte) error {
	p.Status = payload[0]
	err := devMgmtStatusCheck(p.Status)
	if err != nil {
		return err
	}
	bytes := payload[1:]
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

type SetRTCAlarmReq struct {
	wimodMessageImpl
	AlarmType byte
	Hour      byte
	Minutes   byte
	Seconds   byte
}

const (
	AlarmSingle = byte(0x00)
	AlarmDaily  = byte(0x01)
)

func NewSetRTCAlarmReq(alarmType, hour, minutes, seconds byte) *SetRTCAlarmReq {
	req := &SetRTCAlarmReq{}
	req.Init()
	req.AlarmType = alarmType
	req.Hour = hour
	req.Minutes = minutes
	req.Seconds = seconds
	return req
}

func (p *SetRTCAlarmReq) Init() {
	p.code = DEVMGMT_MSG_SET_RTC_ALARM_REQ
}

func (p *SetRTCAlarmReq) String() string {
	return fmt.Sprintf("SetRTCAlarmReq[Type: %X, Hour: %d, Minutes: %d, Seconds: %d]", p.AlarmType, p.Hour, p.Minutes, p.Seconds)
}

func (p *SetRTCAlarmReq) Encode() ([]byte, error) {
	buff := make([]byte, 4)
	buff[0] = p.AlarmType
	buff[1] = p.Hour
	buff[2] = p.Minutes
	buff[3] = p.Seconds
	return buff, nil
}

// DEVMGMT_MSG_SET_RTC_ALARM_RSP

type SetRTCAlarmResp struct {
	wimodMessageStatusImpl
}

func NewSetRTCAlarmResp() *SetRTCAlarmResp {
	resp := &SetRTCAlarmResp{}
	resp.Init()
	return resp
}

func (p *SetRTCAlarmResp) Init() {
	p.code = DEVMGMT_MSG_SET_RTC_ALARM_RSP
}

func (p *SetRTCAlarmResp) String() string {
	return fmt.Sprintf("SetRTCAlarmResp[]")
}

func (p *SetRTCAlarmResp) Decode(payload []byte) error {
	p.Status = payload[0]
	return devMgmtStatusCheck(p.Status)
}

// DEVMGMT_MSG_CLEAR_RTC_ALARM_REQ

type ClearRTCAlarmReq struct {
	wimodMessageImpl
}

func NewClearRTCAlarmReq() *ClearRTCAlarmReq {
	req := &ClearRTCAlarmReq{}
	req.Init()
	return req
}

func (p *ClearRTCAlarmReq) Init() {
	p.code = DEVMGMT_MSG_CLEAR_RTC_ALARM_REQ
}

func (p *ClearRTCAlarmReq) String() string {
	return fmt.Sprintf("ClearRTCAlarmReq[]")
}

func (p *ClearRTCAlarmReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_CLEAR_RTC_ALARM_RSP

type ClearRTCAlarmResp struct {
	wimodMessageStatusImpl
}

func NewClearRTCAlarmResp() *ClearRTCAlarmResp {
	resp := &ClearRTCAlarmResp{}
	resp.Init()
	return resp
}

func (p *ClearRTCAlarmResp) Init() {
	p.code = DEVMGMT_MSG_CLEAR_RTC_ALARM_RSP
}

func (p *ClearRTCAlarmResp) String() string {
	return fmt.Sprintf("ClearRTCAlarmResp[]")
}

func (p *ClearRTCAlarmResp) Decode(payload []byte) error {
	p.Status = payload[0]
	return devMgmtStatusCheck(p.Status)
}

// DEVMGMT_MSG_GET_RTC_ALARM_REQ

type GetRTCAlarmReq struct {
	wimodMessageImpl
}

func NewGetRTCAlarmReq() *GetRTCAlarmReq {
	req := &GetRTCAlarmReq{}
	req.Init()
	return req
}

func (p *GetRTCAlarmReq) Init() {
	p.code = DEVMGMT_MSG_GET_RTC_ALARM_REQ
}

func (p *GetRTCAlarmReq) String() string {
	return fmt.Sprintf("GetRTCAlarmReq[]")
}

func (p *GetRTCAlarmReq) Encode() ([]byte, error) {
	return []byte{}, nil
}

// DEVMGMT_MSG_GET_RTC_ALARM_RSP

type GetRTCAlarmResp struct {
	wimodMessageStatusImpl
	AlarmStatus byte
	AlarmType   byte
	Hour        byte
	Minutes     byte
	Seconds     byte
}

func NewGetRTCAlarmResp() *GetRTCAlarmResp {
	resp := &GetRTCAlarmResp{}
	resp.Init()
	return resp
}

func (p *GetRTCAlarmResp) Init() {
	p.code = DEVMGMT_MSG_GET_RTC_ALARM_RSP
}

func (p *GetRTCAlarmResp) String() string {
	return fmt.Sprintf("GetRTCAlarmResp[Status: %X, Type: %X, Hour: %d, Minutes: %d, Seconds: %d]", p.AlarmStatus, p.AlarmType, p.Hour, p.Minutes, p.Seconds)
}

func (p *GetRTCAlarmResp) Decode(payload []byte) error {
	p.Status = payload[0]
	err := devMgmtStatusCheck(p.Status)
	if err != nil {
		return err
	}
	bytes := payload[1:]
	p.AlarmStatus = bytes[0]
	p.AlarmType = bytes[1]
	p.Hour = bytes[2]
	p.Minutes = bytes[3]
	p.Seconds = bytes[4]
	return nil
}

// DEVMGMT_MSG_RTC_ALARM_IND

type RTCAlarmInd struct {
	wimodMessageStatusImpl
}

func NewRTCAlarmInd() *RTCAlarmInd {
	ind := &RTCAlarmInd{}
	ind.Init()
	return ind
}

func (p *RTCAlarmInd) Init() {
	p.code = DEVMGMT_MSG_RTC_ALARM_IND
}

func (p *RTCAlarmInd) String() string {
	return fmt.Sprintf("RTCAlarmInd[Status: 0x%02X]", p.Status)
}

func (p *RTCAlarmInd) Decode(payload []byte) error {
	p.Status = payload[0]
	return devMgmtStatusCheck(p.Status)
}
