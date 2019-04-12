package slip

import (
	"bytes"

	"github.com/tarm/serial"
)

const (
	SLIP_END     = 0300 /* 0xC0 indicates end of packet */
	SLIP_ESC     = 0333 /* 0xDB, indicates byte stuffing */
	SLIP_ESC_END = 0334 /* 0xDC, ESC ESC_END means END data byte */
	SLIP_ESC_ESC = 0335 /* 0xDD, ESC ESC_ESC means ESC data byte */
)

const (
	BUFFER_SIZE = 128
)

const (
	SLIPDEC_STATE_START = iota
	SLIPDEC_STATE_IN_FRAME
	SLIPDEC_STATE_ESC
)

type SlipPacket struct {
	Buffer bytes.Buffer
	Err    error
}

func SlipDecoder(s *serial.Port, c chan<- SlipPacket) {
	buf := make([]byte, BUFFER_SIZE)
	slipPacket := SlipPacket{}
	slipPacket.Buffer.Grow(BUFFER_SIZE)
	state := SLIPDEC_STATE_START
	for {
		n, err := s.Read(buf)
		if err != nil {
			errInf := SlipPacket{Err: err}
			c <- errInf
		}
		for _, b := range buf[:n] {
			switch state {
			case SLIPDEC_STATE_START:
				if b == SLIP_END {
					state = SLIPDEC_STATE_IN_FRAME
				}
				slipPacket.Buffer.Reset()
				break
			case SLIPDEC_STATE_IN_FRAME:
				switch b {
				case SLIP_END:
					if slipPacket.Buffer.Len() != 0 {
						c <- slipPacket
					}
					slipPacket.Buffer.Reset()
					state = SLIPDEC_STATE_START
					break
				case SLIP_ESC:
					state = SLIPDEC_STATE_ESC
					break
				default:
					slipPacket.Buffer.WriteByte(b)
					break
				}
				break
			case SLIPDEC_STATE_ESC:
				switch b {
				case SLIP_ESC_END:
					slipPacket.Buffer.WriteByte(SLIP_END)
					state = SLIPDEC_STATE_IN_FRAME
					break
				case SLIP_ESC_ESC:
					slipPacket.Buffer.WriteByte(SLIP_ESC)
					state = SLIPDEC_STATE_IN_FRAME
					break
				default:
					state = SLIPDEC_STATE_START
					break
				}
				break
			}
		}
	}
}
