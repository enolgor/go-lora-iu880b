package iu880b

import (
	"fmt"
	"io"
	"sync"

	"./hci"
	"./slip"
	"./wimod"
)

type WiModController struct {
	rwc           io.ReadWriteCloser
	slipDecoder   *slip.SlipDecoder
	closer        chan bool
	events        chan wimod.WiModMessageInd
	respChannels  map[uint16][]chan hci.HCIPacket
	eventChannels []chan wimod.WiModMessageInd
	noBlock       bool
	mutex         *sync.Mutex
}

type WiModControllerConfig struct {
	Stream          io.ReadWriteCloser
	EventBufferSize int
	EventNoBlock    bool
}

func NewController(config *WiModControllerConfig) *WiModController {
	slipDecoder := slip.NewDecoder(config.Stream)
	respChannels := make(map[uint16][]chan hci.HCIPacket)
	eventChannels := []chan wimod.WiModMessageInd{}
	eventBufferSize := 10
	if config.EventBufferSize != 0 {
		eventBufferSize = config.EventBufferSize
	}
	events := make(chan wimod.WiModMessageInd, eventBufferSize)
	closer := make(chan bool, 1)
	controller := &WiModController{config.Stream, &slipDecoder, closer, events, respChannels, eventChannels, config.EventNoBlock, &sync.Mutex{}}
	go controller.start()
	go controller.eventDispatcher()
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
			if wimod.IsAlarm(code) {
				ind, err := wimod.DecodeInd(&hciPacket)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					continue
				} else {
					if c.noBlock && len(c.events) == cap(c.events) {
						discarded := <-c.events
						fmt.Printf("Event buffer full. Discarding oldest event: [%04X]\n", discarded.Code())
					}
					c.events <- ind
					continue
				}
			}
			c.mutex.Lock()
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
			c.mutex.Unlock()
		}
	}
}

func (c *WiModController) eventDispatcher() {
	for event := range c.events {
		c.mutex.Lock()
		channels := c.eventChannels
		c.eventChannels = nil
		for _, channel := range channels {
			channel <- event
			close(channel)
		}
		c.mutex.Unlock()
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
	c.mutex.Lock()
	channels, ok := c.respChannels[resp.Code()]
	if !ok {
		channels = []chan hci.HCIPacket{}
	}
	channels = append(channels, respChannel)
	c.respChannels[resp.Code()] = channels
	c.mutex.Unlock()

	err := c.sendReq(req)
	if err != nil {
		return err
	}
	hci := <-respChannel
	err = wimod.DecodeResp(&hci, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c *WiModController) ReadInd() wimod.WiModMessageInd {
	eventChannel := make(chan wimod.WiModMessageInd)
	c.mutex.Lock()
	c.eventChannels = append(c.eventChannels, eventChannel)
	c.mutex.Unlock()
	event := <-eventChannel
	return event
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
