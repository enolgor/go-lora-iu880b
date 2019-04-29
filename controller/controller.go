package controller

import (
	"fmt"
	"io"
	"sync"

	"github.com/enolgor/wimod-lorawan-endnode-controller/hci"
	"github.com/enolgor/wimod-lorawan-endnode-controller/slip"
	"github.com/enolgor/wimod-lorawan-endnode-controller/wimod"
)

type WiModController struct {
	rwc           io.ReadWriteCloser
	slipDecoder   *slip.SlipDecoder
	closer        chan bool
	events        chan hci.HCIPacket
	respChannels  map[uint16][]chan hci.HCIPacket
	eventChannels map[uint16][]chan hci.HCIPacket
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
	eventChannels := make(map[uint16][]chan hci.HCIPacket)
	eventBufferSize := 10
	if config.EventBufferSize != 0 {
		eventBufferSize = config.EventBufferSize
	}
	events := make(chan hci.HCIPacket, eventBufferSize)
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
				if c.noBlock && len(c.events) == cap(c.events) {
					discarded := <-c.events
					fmt.Printf("Event buffer full. Discarding oldest event: [%02X%02X]\n", discarded.Dst, discarded.ID)
				}
				c.events <- hciPacket
				continue
			}
			c.mutex.Lock()
			channels, ok := c.respChannels[code]
			if !ok {
				fmt.Printf("Discarded packet because no listener: %v\n", hciPacket)
				c.mutex.Unlock()
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
		code := (uint16(event.Dst) << 8) + uint16(event.ID)
		c.mutex.Lock()
		channels := c.eventChannels[code]
		channelsAll := c.eventChannels[0]
		c.eventChannels[0] = nil
		c.eventChannels[code] = nil
		channels = append(channels, channelsAll...)
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
	req.Init()
	resp.Init()
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
	return wimod.DecodeResp(&hci, resp)
}

func (c *WiModController) ReadInd() (wimod.WiModMessageInd, error) {
	eventChannel := make(chan hci.HCIPacket)
	c.mutex.Lock()
	channels, ok := c.eventChannels[0]
	if !ok {
		channels = []chan hci.HCIPacket{}
	}
	channels = append(channels, eventChannel)
	c.eventChannels[0] = channels
	c.mutex.Unlock()
	hci := <-eventChannel
	return wimod.DecodeInd(&hci)
}

func (c *WiModController) ReadSpecificInd(ind wimod.WiModMessageInd) error {
	ind.Init()
	eventChannel := make(chan hci.HCIPacket)
	c.mutex.Lock()
	channels, ok := c.eventChannels[ind.Code()]
	if !ok {
		channels = []chan hci.HCIPacket{}
	}
	channels = append(channels, eventChannel)
	c.eventChannels[ind.Code()] = channels
	c.mutex.Unlock()
	hci := <-eventChannel
	return wimod.DecodeSpecificInd(&hci, ind)
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
