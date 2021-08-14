package cloudapp

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/cws"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/webrtc"

	"github.com/pion/rtp"
)

const (
	DefaultSTUNTURN = `[{"urls":"stun:stun.l.google.com:19302"}]`
)

var appEventTypes []string = []string{"MOUSEDOWN", "MOUSEUP", "MOUSEMOVE", "KEYDOWN", "KEYUP"}

type Service struct {
	clients map[string]*Client
	hosts map[string]*Client
	ccApp   CloudAppClient
	config  config.Config
	// communicate with cloud app
	appEvents chan Packet
}

type Client struct {
	clientID    string
	ws          *cws.Client
	rtcConn     *webrtc.WebRTC
	videoStream chan *rtp.Packet
	audioStream chan *rtp.Packet
	appEvents   chan Packet
	// cancel to trigger cleaning up when client is disconnected
	cancel chan struct{}
	// done to notify if the client is done clean up
	done chan struct{}
	// TODO: Get rid of ssrc
	ssrc uint32
}

type AppHost struct {
	// Host string `json:"host"`
	Addr    string `json:"addr"`
	AppName string `json:"app_name"`
}

func (s *Service) AddClient(clientID string, ws *cws.Client) *Client {
	client := NewServiceClient(clientID, ws, s.appEvents, s.ccApp.GetSSRC())
	s.clients[clientID] = client
	return client
}

func (s *Service) RemoveClient(clientID string) {
	client := s.clients[clientID]
	close(client.cancel)
	<-client.done
	if client.rtcConn != nil {
		client.rtcConn.StopClient()
		client.rtcConn = nil
	}
}

func NewServiceClient(clientID string, ws *cws.Client, appEvents chan Packet, ssrc uint32) *Client {
	// The 1st packet
	// Note: We won't force browser to initwebrtc yet when new host connects
	//ws.Send(cws.WSPacket{
	//	Type: "init",
	//	Data: DefaultSTUNTURN,
	//}, nil)

	return &Client{
		appEvents:   appEvents,
		clientID:    clientID,
		ws:          ws,
		ssrc:        ssrc,
		videoStream: make(chan *rtp.Packet, 100),
		audioStream: make(chan *rtp.Packet, 100),
		cancel:      make(chan struct{}),
		done:        make(chan struct{}),
	}
}

func (c *Client) Handle() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered when sent to close Image Channel")
		}
	}()

	wg := sync.WaitGroup{}

	// Video Stream
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered. Maybe we :sent to Closed Channel", r)
				wg.Done()
			}
		}()

	loop:
		for packet := range c.videoStream {
			select {
			case <-c.cancel:
				break loop
			case c.rtcConn.ImageChannel <- packet:
			}
		}
		wg.Done()
		log.Println("Closed Service Video Channel")
	}()

	// Audio Stream
	// wg.Add(1)
	// go func() {
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			log.Println("Recovered. Maybe we :sent to Closed Channel", r)
	// 			wg.Done()
	// 		}
	// 	}()

	// loop:
	// 	for packet := range c.audioStream {
	// 		select {
	// 		case <-c.cancel:
	// 			break loop
	// 		case c.rtcConn.AudioChannel <- packet:
	// 		}
	// 	}
	// 	wg.Done()
	// 	log.Println("Closed Service Audio Channel")
	// }()

	// Input stream is closed after StopClient . TODO: check if can close earlier
	// wg.Add(1)
	go func() {
		// Data channel input
		for rawInput := range c.rtcConn.InputChannel {
			// TODO: No dynamic allocation
			wspacket := cws.WSPacket{}
			err := json.Unmarshal(rawInput, &wspacket)
			if err != nil {
				log.Println(err)
			}
			c.appEvents <- convertWSPacket(wspacket)
		}
	}()
	wg.Wait()
	close(c.done)
}

// Route: Handshake to initialize WebRTC
func (c *Client) Route(ssrc uint32, s *Server) {

	//ws.Send(cws.WSPacket{
	//	Type: "init",
	//	Data: DefaultSTUNTURN,
	//}, nil)

	c.ws.Receive(
		"init",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			var appPath = req.Data
			s.capp.ccApp = NewCloudAppClient(s.capp.config, s.capp.appEvents, appPath)
			c.ws.Send(cws.WSPacket{
				Type: "init",
				Data: DefaultSTUNTURN,
			}, nil)
			return cws.EmptyPacket
		},
	)

	c.ws.Receive(
		"INIT",
		func(resp cws.WSPacket) (req cws.WSPacket) {
			log.Println("Received INIT")

			return cws.EmptyPacket
		},
	)

	c.ws.Receive("initwebrtc", func(req cws.WSPacket) (resp cws.WSPacket) {
		log.Println("Received a request to createOffer from browser", req)

		c.rtcConn = webrtc.NewWebRTC()
		var initPacket struct {
			IsMobile bool `json:"is_mobile"`
		}
		err := json.Unmarshal([]byte(req.Data), &initPacket)
		if err != nil {
			log.Println("Error: Cannot decode json:", err)
			return cws.EmptyPacket
		}

		localSession, err := c.rtcConn.StartClient(
			initPacket.IsMobile,
			func(candidate string) {
				// send back candidate string to browser
				c.ws.Send(cws.WSPacket{
					Type:      "candidate",
					Data:      candidate,
					SessionID: req.SessionID,
				}, nil)
			},
			ssrc,
		)

		if err != nil {
			log.Println("Error: Cannot create new webrtc session", err)
			return cws.EmptyPacket
		}

		return cws.WSPacket{
			Type: "offer",
			Data: localSession,
		}
	})

	c.ws.Receive(
		"answer",
		func(resp cws.WSPacket) (req cws.WSPacket) {
			log.Println("Received answer SDP from browser", resp)
			err := c.rtcConn.SetRemoteSDP(resp.Data)
			if err != nil {
				log.Println("Error: Cannot set RemoteSDP of client: " + resp.SessionID)
			}

			go c.Handle()
			return cws.EmptyPacket
		},
	)

	c.ws.Receive(
		"candidate",
		func(resp cws.WSPacket) (req cws.WSPacket) {
			log.Println("Received remote Ice Candidate from browser")

			err := c.rtcConn.AddCandidate(resp.Data)
			if err != nil {
				log.Println("Error: Cannot add IceCandidate of client: " + resp.SessionID)
			}

			return cws.EmptyPacket
		},
	)

	c.ws.Receive(
		"heartbeat",
		func(resp cws.WSPacket) (req cws.WSPacket) {
			log.Println("Received heartbeat")
			return cws.EmptyPacket
		},
	)
}

// NewCloudService returns a Cloud Service
func NewCloudService(cfg config.Config) *Service {
	appEvents := make(chan Packet, 1)
	s := &Service{
		clients:   map[string]*Client{},
		appEvents: appEvents,
		// ccApp is only initiated later in "init" websocket message
		ccApp:     NewCloudAppClient(cfg, appEvents, ""),
		config:    cfg,
	}

	return s
}

// func (s *Service) SendInput(packet Packet) {
// 	s.ccApp.SendInput(packet)
// }

func (s *Service) GetSSRC() uint32 {
	return s.ccApp.GetSSRC()
}

func (s *Service) Handle() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered when sent to closed Video Stream channel", r)
			}
		}()
		for p := range s.ccApp.VideoStream() {
			for id, client := range s.clients {
				select {
				case <-client.cancel:
					log.Println("Closing Video Audio")
					// stop producing for client
					delete(s.clients, id)
					// close(client.audioStream)
					close(client.videoStream)
				case client.videoStream <- p:
				}
			}
		}
	}()
	// go func() {
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			log.Println("Recovered when sent to closed Video Stream channel", r)
	// 		}
	// 	}()
	// 	for p := range s.ccApp.AudioStream() {
	// 		for _, client := range s.clients {
	// 			select {
	// 			// case <-client.cancel:
	// 			// fmt.Println("Closing Audio")
	// 			// stop producing for client
	// 			// close(client.audioStream)
	// 			case client.audioStream <- p:
	// 			}
	// 		}
	// 	}
	// }()
	//s.ccApp.Handle()
}
