package wimod

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

func EncodeRTCTime(t time.Time) uint32 {
	month := int(t.Month())
	day := t.Day()
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()
	year := t.Year() - 2000
	rtc := uint32(0)
	rtc += uint32((year & 0x3F) << 26)
	rtc += uint32((day & 0x1F) << 21)
	rtc += uint32((hour & 0x1F) << 16)
	rtc += uint32((month & 0x0F) << 12)
	rtc += uint32((minute & 0x3F) << 6)
	rtc += uint32((second & 0x3F))
	return rtc
}

func DecodeRTCTime(rtc uint32) time.Time {
	month := int((rtc >> 12) & 0x0F)
	day := int((rtc >> 21) & 0x1F)
	hour := int((rtc >> 16) & 0x1F)
	minute := int((rtc >> 6) & 0x3F)
	second := int(rtc & 0x3F)
	year := 2000 + int((rtc>>26)&0x3F)
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
}

type EUI uint64

func EncodeEUI(eui *EUI) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(*eui))
	return bytes
}

func DecodeEUI(bytes []byte) EUI {
	euibytes := make([]byte, 8)
	euibytes = append(euibytes[:8-len(bytes)], bytes...)
	return EUI(binary.BigEndian.Uint64(euibytes))
}

func (e EUI) String() string {
	return fmt.Sprintf("%016X", uint64(e))
}

func ParseEUI(str string) (EUI, error) {
	v := strconv.ParseUint(str, 16, 64)
	return EUI(v)
}

type Key [2]uint64

func EncodeKey(key *Key) []byte {
	bytes := make([]byte, 16)
	binary.BigEndian.PutUint64(bytes[:8], uint64(key[0]))
	binary.BigEndian.PutUint64(bytes[8:], uint64(key[1]))
	return bytes
}

func DecodeKey(bytes []byte) Key {
	keybytes := make([]byte, 16)
	keybytes = append(keybytes[:16-len(bytes)], bytes...)
	key := [2]uint64{}
	key[0] = binary.BigEndian.Uint64(keybytes[:8])
	key[1] = binary.BigEndian.Uint64(keybytes[8:])
	return Key(key)
}

func (k Key) String() string {
	return fmt.Sprintf("%016X%016X", k[0], k[1])
}
