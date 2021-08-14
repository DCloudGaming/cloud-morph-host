package cloudapp

import (
	"encoding/json"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/cws"
	"log"
)

const (
	DefaultSTUNTURN = `[{"urls":"stun:stun.l.google.com:19302"}]`
)

type ChosenHostApp struct {
	hostID string
	appPath string
}

type Service struct {
	clients map[string]*Client
	hosts   map[string]*Host
	ccApp   *ccImpl
	config  config.Config
	clientChosenApp map[string]ChosenHostApp
}

type Client struct {
	clientID string
	ws       *cws.Client
	// cancel to trigger cleaning up when client is disconnected
	cancel chan struct{}
	// done to notify if the client is done clean up
	done chan struct{}
}

type Host struct {
	hostID string
	ws     *cws.Client
	appPaths []string
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
		clientID: clientID,
		ws:       ws,
		cancel:   make(chan struct{}),
		done:     make(chan struct{}),
	}
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
		hostID: hostID,
		ws:     ws,
		cancel: make(chan struct{}),
		done:   make(chan struct{}),
	}
}

func addForwardingRoute(sender *cws.Client, senderID string, receiver *cws.Client, receiverID string, messages []string, s *Server, is_sender_browser bool) {
	for _, message := range messages {
		sender.Receive(
			message,
			func(req cws.WSPacket) cws.WSPacket {

				if (is_sender_browser) {
					if (s.capp.clientChosenApp[senderID].hostID != receiverID) {
						return cws.EmptyPacket
					}
				} else {
					if (s.capp.clientChosenApp[receiverID].hostID != senderID) {
						return cws.EmptyPacket
					}
				}

				resp := receiver.SyncSend(req)
				return resp
			},
		)
	}
}

func (h *Host) HostRoute(s *Server) {
	h.ws.Receive(
		"registerApps",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			var registerAppsData struct {
				AppPaths []string `json:"app_paths"`
			}
			log.Println("Get app registrations from host")
			err := json.Unmarshal([]byte(req.Data), &registerAppsData)
			if err != nil {
				log.Println("Error: Cannot decode json:" , err)
				return cws.EmptyPacket
			}

			h.appPaths = registerAppsData.AppPaths

			type MinimalHostMeta struct {
				hostID string `json:"host_id"`
				appPaths []string `json:"app_paths"`
			}

			type HostsMeta struct {
				hosts []MinimalHostMeta	`json:"hosts"`
			}

			var hosts []MinimalHostMeta
			for _, h := range s.capp.hosts {
				hosts = append(hosts, MinimalHostMeta{h.hostID, h.appPaths})
			}

			hostsData := HostsMeta{
				hosts: hosts,
			}

			hostsJsonData, err := json.Marshal(hostsData)
			if err != nil {
				return
			}

			for _, client := range s.capp.clients {
				client.ws.Send(cws.WSPacket{
					Type: "hostsUpdated",
					Data: string(hostsJsonData),
				}, nil)
			}

			return cws.EmptyPacket
		},
		)
}

func (c *Client) ClientRoute(s *Server) {
	c.ws.Receive(
			"registerBrowserHost",
			func(req cws.WSPacket) (resp cws.WSPacket) {
				var hostAppsData struct {
					host_id string
					app     string
				}
				err := json.Unmarshal([]byte(req.Data), &hostAppsData)
				if err != nil {
					log.Println("Error: Cannot decode json:" , err)
					return cws.EmptyPacket
				}
				s.capp.clientChosenApp[c.clientID] = ChosenHostApp{
					hostID:  hostAppsData.host_id,
					appPath: hostAppsData.app,
				}

				// Send this request back to host. Host will start ffmpeg,
				// and send the request back to browser to start initwebrtc
				s.capp.hosts[hostAppsData.host_id].ws.Send(cws.WSPacket{
					Type: "init",
					Data: hostAppsData.app,
				}, nil)
				return cws.EmptyPacket
			},
		)
}

// NewCloudService returns a Cloud Service
func NewCloudService(cfg config.Config) *Service {
	s := &Service{
		clients: map[string]*Client{},
		hosts:   map[string]*Host{},
		ccApp:   NewCloudAppClient(cfg),
		config:  cfg,
	}

	return s
}
