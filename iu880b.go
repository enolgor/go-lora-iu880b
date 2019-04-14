package iu880b

import (
	"fmt"
	"io"
	"log"
	"time"

	"./hci"
	"./slip"
	"./wimod"
	"github.com/tarm/serial"
)

func performTest() {
	c := &serial.Config{Name: "COM4", Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	controller := NewController(s)
	req := wimod.NewGetRTCReq()
	resp := wimod.NewGetRTCResp()
	err = controller.Request(req, resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)

	time.Sleep(1 * time.Second)

	req2 := wimod.NewGetRTCReq()
	resp2 := wimod.NewGetRTCResp()
	err = controller.Request(req2, resp2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp2)
	controller.Close()
}

func performTest2() {
	c := &serial.Config{Name: "COM4", Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	sendWakeUp(s)
	req := wimod.NewGetRTCReq()
	sendReq(s, req)

	slipDecoder := slip.NewDecoder(s)
	hciDecoder := hci.NewDecoder(&slipDecoder)
	for {
		hciPacket := hci.HCIPacket{}
		err := hciDecoder.Read(&hciPacket)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		infoResp := wimod.NewGetRTCResp()
		err = wimod.DecodeResp(&hciPacket, infoResp)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		// infoResp, ok := resp.(*wimod.GetRTCResp)
		fmt.Println(infoResp)
	}
}

func sendWakeUp(rw io.ReadWriter) {
	for i := 0; i < 40; i++ {
		rw.Write([]byte{slip.SLIP_END})
	}
}

func sendReq(s *serial.Port, req wimod.WiModMessageReq) error {
	hci, err := wimod.EncodeReq(req)
	if err != nil {
		return err
	}
	slipPacket := slip.SlipEncode(hci.Encode())
	_, err = s.Write(slipPacket)
	return err
}

type WiModController struct {
	rw           io.ReadWriter
	slipDecoder  *slip.SlipDecoder
	hciDecoder   *hci.HciDecoder
	closer       chan bool
	respChannels map[uint16][]chan hci.HCIPacket
}

func NewController(rw io.ReadWriter) *WiModController {
	slipDecoder := slip.NewDecoder(rw)
	hciDecoder := hci.NewDecoder(&slipDecoder)
	respChannels := make(map[uint16][]chan hci.HCIPacket)
	closer := make(chan bool, 1)
	controller := &WiModController{rw, &slipDecoder, &hciDecoder, closer, respChannels}
	go controller.start()
	return controller
}

func (c *WiModController) start() {
	hciPacket := hci.HCIPacket{}
	for {
		select {
		case <-c.closer:
			fmt.Println("closing!")
			return
		default:
			err := c.hciDecoder.Read(&hciPacket)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				continue
			}
			respCode := (uint16(hciPacket.Dst) << 8) + uint16(hciPacket.ID)
			channels, ok := c.respChannels[respCode]
			if !ok {
				fmt.Printf("Discarded packet because no listener: %v", hciPacket)
				continue
			}
			respChannel := channels[0]
			respChannel <- hciPacket
			close(respChannel)
			channels[0] = nil
			channels = channels[1:]
			c.respChannels[respCode] = channels
		}
	}
}

func (c *WiModController) Close() {
	c.closer <- true
	pingReq := wimod.NewPingReq()
	pingResp := wimod.NewPingResp()
	c.Request(pingReq, pingResp)
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

func (c *WiModController) sendReq(req wimod.WiModMessageReq) error {
	hci, err := wimod.EncodeReq(req)
	if err != nil {
		return err
	}
	slipPacket := slip.SlipEncode(hci.Encode())
	sendWakeUp(c.rw)
	_, err = c.rw.Write(slipPacket)
	return err
}
