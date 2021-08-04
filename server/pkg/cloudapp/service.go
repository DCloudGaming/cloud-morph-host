package cloudapp

import (
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/cws"
)

type Service struct {
	clients map[string]*Client
	hosts map[string]*Host
	ccApp   *ccImpl
	config  config.Config
}

type Client struct {
	clientID    string
	ws          *cws.Client
	// cancel to trigger cleaning up when client is disconnected
	cancel chan struct{}
	// done to notify if the client is done clean up
	done chan struct{}
}

type Host struct {
	hostID    string
	ws          *cws.Client
	// cancel to trigger cleaning up when client is disconnected
	cancel chan struct{}
	// done to notify if the client is done clean up
	done chan struct{}
}

type AppHost struct {
	// Host string `json:"host"`
	Addr    string `json:"addr"`
	AppName string `json:"app_name"`
}

func (s *Service) AddClient(clientID string, ws *cws.Client) *Client {
	client := NewServiceClient(clientID, ws)
	s.clients[clientID] = client
	return client
}

func (s *Service) RemoveClient(clientID string) {
	client := s.clients[clientID]
	close(client.cancel)
}

func NewServiceClient(clientID string, ws *cws.Client) *Client {
	return &Client{
		clientID:    clientID,
		ws:          ws,
		cancel:      make(chan struct{}),
		done:        make(chan struct{}),
	}
}

// Route: Handshake to initialize WebRTC
func (c *Client) Route(hosts map[string]*Host) {
	c.ws.Receive("initwebrtc",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			for _, host := range hosts {
				host.ws.Send(req, nil)
			}
			return cws.EmptyPacket
		})

	c.ws.Receive(
		"answer",
		func(resp cws.WSPacket) (req cws.WSPacket) {
			for _, host := range hosts {
				host.ws.Send(req, nil)
			}
			return cws.EmptyPacket
		})

	c.ws.Receive(
		"candidate",
		func(resp cws.WSPacket) (req cws.WSPacket) {
			for _, host := range hosts {
				host.ws.Send(req, nil)
			}
			return cws.EmptyPacket
		})
}

func (s *Service) AddHost(hostID string, ws *cws.Client) *Host {
	host := NewServiceHost(hostID, ws)
	s.hosts[hostID] = host
	return host
}

func (s *Service) RemoveHost(hostID string) {
	host := s.hosts[hostID]
	close(host.cancel)
}

func NewServiceHost(hostID string, ws *cws.Client) *Host {
	return &Host{
		hostID:    hostID,
		ws:          ws,
		cancel:      make(chan struct{}),
		done:        make(chan struct{}),
	}
}

// Route: Handshake to initialize WebRTC
func (c *Host) Route(clients map[string]*Client) {
	c.ws.Receive(
		"init",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			for _, client := range clients {
				client.ws.Send(req, nil)
			}
			return cws.EmptyPacket
		})

	c.ws.Receive(
		"INIT",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			for _, client := range clients {
				client.ws.Send(req, nil)
			}
			return cws.EmptyPacket
		})

	c.ws.Receive(
		"candidate",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			for _, client := range clients {
				client.ws.Send(req, nil)
			}
			return cws.EmptyPacket
		})

	c.ws.Receive(
		"offer",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			for _, client := range clients {
				client.ws.Send(req, nil)
			}
			return cws.EmptyPacket
		})
}

// NewCloudService returns a Cloud Service
func NewCloudService(cfg config.Config) *Service {
	s := &Service{
		clients:   map[string]*Client{},
		ccApp:     NewCloudAppClient(cfg),
		config:    cfg,
	}

	return s
}
