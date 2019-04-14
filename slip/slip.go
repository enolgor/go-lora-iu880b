package slip

import (
	"bytes"
	"io"
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

type slipPacket struct {
	Buffer bytes.Buffer
	Err    error
}

func SlipEncode(payload []byte) []byte {
	var buff bytes.Buffer
	buff.WriteByte(SLIP_END)
	for _, b := range payload {
		switch b {
		case SLIP_END:
			buff.WriteByte(SLIP_ESC)
			buff.WriteByte(SLIP_ESC_END)
			break
		case SLIP_ESC:
			buff.WriteByte(SLIP_ESC)
			buff.WriteByte(SLIP_ESC_ESC)
			break
		default:
			buff.WriteByte(b)
			break
		}
	}
	buff.WriteByte(SLIP_END)
	return buff.Bytes()
}

type SlipDecoder struct {
	reader io.Reader
	ch     chan slipPacket
}

func (d *SlipDecoder) Read() ([]byte, error) {
	packet := <-d.ch
	if packet.Err != nil {
		return nil, packet.Err
	}
	return packet.Buffer.Bytes(), nil
}

/*
func (d *SlipDecoder) Close() {
	close(d.ch)
}*/

func NewDecoder(reader io.Reader) SlipDecoder {
	ch := make(chan slipPacket)
	go slipChannelDecoder(reader, ch)
	return SlipDecoder{reader, ch}
}

func slipChannelDecoder(s io.Reader, c chan<- slipPacket) {
	buf := make([]byte, BUFFER_SIZE)
	packet := slipPacket{}
	packet.Buffer.Grow(BUFFER_SIZE)
	state := SLIPDEC_STATE_START
	for {
		n, err := s.Read(buf)
		if err != nil {
			errInf := slipPacket{Err: err}
			c <- errInf
		}
		for _, b := range buf[:n] {
			switch state {
			case SLIPDEC_STATE_START:
				if b == SLIP_END {
					state = SLIPDEC_STATE_IN_FRAME
				}
				packet.Buffer.Reset()
				break
			case SLIPDEC_STATE_IN_FRAME:
				switch b {
				case SLIP_END:
					if packet.Buffer.Len() != 0 {
						c <- packet
					}
					packet.Buffer.Reset()
					state = SLIPDEC_STATE_START
					break
				case SLIP_ESC:
					state = SLIPDEC_STATE_ESC
					break
				default:
					packet.Buffer.WriteByte(b)
					break
				}
				break
			case SLIPDEC_STATE_ESC:
				switch b {
				case SLIP_ESC_END:
					packet.Buffer.WriteByte(SLIP_END)
					state = SLIPDEC_STATE_IN_FRAME
					break
				case SLIP_ESC_ESC:
					packet.Buffer.WriteByte(SLIP_ESC)
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
