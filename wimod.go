package iu880b

import (
	"encoding/binary"
	"fmt"
)

type WiModDevStatus struct {
	Status               byte
	SystemTickResolution byte
	SystemTicks          uint32
	TargetTime           RTCTime
	NVMStatus            uint16
	BatteryLevel         uint16
}

func (w *WiModDevStatus) Decode(payload []byte) {
	w.Status = payload[0]
	w.SystemTickResolution = payload[1]
	w.SystemTicks = binary.LittleEndian.Uint32(payload[2:6])
	w.TargetTime = RTCTime(binary.LittleEndian.Uint32(payload[6:10]))
	w.NVMStatus = binary.LittleEndian.Uint16(payload[10:12])
	w.BatteryLevel = binary.LittleEndian.Uint16(payload[12:14])
}

func (w WiModDevStatus) String() string {
	return fmt.Sprintf("Battery Level: %d, Time: %s", w.BatteryLevel, w.TargetTime)
}

type RTCTime uint32

func (r RTCTime) String() string {
	return fmt.Sprintf("%d-%d-%d %d:%d:%d", r.getYears(), r.getMonths(), r.getDays(), r.getHours(), r.getMinutes(), r.getSeconds())
}

func (r RTCTime) getMonths() int {
	return int((uint32(r) >> 12) & 0x0F)
}

func (r RTCTime) getDays() int {
	return int((uint32(r) >> 21) & 0x1F)
}

func (r RTCTime) getHours() int {
	return int((uint32(r) >> 16) & 0x1F)
}

func (r RTCTime) getMinutes() int {
	return int((uint32(r) >> 6) & 0x3F)
}

func (r RTCTime) getSeconds() int {
	return int(uint32(r) & 0x3F)
}

func (r RTCTime) getYears() int {
	return 2000 + int((uint32(r)>>26)&0x3F)
}
