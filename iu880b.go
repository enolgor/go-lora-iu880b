package iu880b

import (
	"fmt"
	"io"

	"./hci"
	"./slip"
	"./wimod"
)

type WiModController struct {
	rwc          io.ReadWriteCloser
	slipDecoder  *slip.SlipDecoder
	closer       chan bool
	events       chan wimod.WiModMessageInd
	respChannels map[uint16][]chan hci.HCIPacket
	noBlock      bool
}

type WiModControllerConfig struct {
	Stream          io.ReadWriteCloser
	EventBufferSize int
	EventNoBlock    bool
}

func NewController(config *WiModControllerConfig) *WiModController {
	slipDecoder := slip.NewDecoder(config.Stream)
	respChannels := make(map[uint16][]chan hci.HCIPacket)
	eventBufferSize := 10
	if config.EventBufferSize != 0 {
		eventBufferSize = config.EventBufferSize
	}
	events := make(chan wimod.WiModMessageInd, eventBufferSize)
	closer := make(chan bool, 1)
	controller := &WiModController{config.Stream, &slipDecoder, closer, events, respChannels, config.EventNoBlock}
	go controller.start()
	return controller
}

func (c *WiModController) start() {
	hciPacket := hci.HCIPacket{}
	run := true
	for run {
		select {
		case <-c.closer:
			run = false
			break
		default:
			payload, err := c.slipDecoder.Read()
			if len(payload) == 0 {
				continue
			}
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				continue
			}
			err = hciPacket.Decode(payload)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				continue
			}
			code := (uint16(hciPacket.Dst) << 8) + uint16(hciPacket.ID)
			if wimod.IsAlarm(hciPacket.Dst, hciPacket.ID) {
				ind, err := wimod.DecodeInd(&hciPacket)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					continue
				} else {
					if c.noBlock && len(c.events) == cap(c.events) {
						discarded := <-c.events
						fmt.Printf("Event buffer full. Discarding oldest event: [%02X][%02X]\n", discarded.GetDst(), discarded.GetID())
					}
					c.events <- ind
					continue
				}
			}
			channels, ok := c.respChannels[code]
			if !ok {
				fmt.Printf("Discarded packet because no listener: %v", hciPacket)
				continue
			}
			respChannel := channels[0]
			respChannel <- hciPacket
			close(respChannel)
			channels[0] = nil
			channels = channels[1:]
			c.respChannels[code] = channels
		}
	}
}

func (c *WiModController) Close() {
	c.closer <- true
	c.slipDecoder.Close()
	c.rwc.Close()
	close(c.closer)
	close(c.events)
}

func (c *WiModController) Request(req wimod.WiModMessageReq, resp wimod.WiModMessageResp) error {
	respChannel := make(chan hci.HCIPacket)
	respCode := (uint16(resp.GetDst()) << 8) + uint16(resp.GetID())
	channels, ok := c.respChannels[respCode]
	if !ok {
		channels = []chan hci.HCIPacket{}
	}
	err := c.sendReq(req)
	if err != nil {
		return err
	}
	channels = append(channels, respChannel)
	c.respChannels[respCode] = channels
	hci := <-respChannel
	err = wimod.DecodeResp(&hci, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c *WiModController) ReadInd() wimod.WiModMessageInd {
	return <-c.events
}

func (c *WiModController) sendReq(req wimod.WiModMessageReq) error {
	hci, err := wimod.EncodeReq(req)
	if err != nil {
		return err
	}
	slipPacket := slip.SlipEncode(hci.Encode())
	sendWakeUp(c.rwc)
	_, err = c.rwc.Write(slipPacket)
	return err
}

func sendWakeUp(rw io.ReadWriter) {
	for i := 0; i < 40; i++ {
		rw.Write([]byte{slip.SLIP_END})
	}
}
