package wimod

import "time"

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
