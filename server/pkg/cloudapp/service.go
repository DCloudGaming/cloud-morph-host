package cloudapp

import (
	"encoding/json"
	"log"

	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/cws"
)

const (
	DefaultSTUNTURN = `[{"urls":"stun:stun.l.google.com:19302"}]`
)

type ChosenHostApp struct {
	hostID  string
	appPath string
}

type Service struct {
	clients         map[string]*Client
	hosts           map[string]*Host
	ccApp           *ccImpl
	config          config.Config
	clientChosenApp map[string]ChosenHostApp
}

type Client struct {
	clientID string
	ws       *cws.Client
	// cancel to trigger cleaning up when client is disconnected
	cancel chan struct{}
	// done to notify if the client is done clean up
	done chan struct{}
	walletAddress string
}

type Host struct {
	hostID string
	walletAddress string
	ws     *cws.Client
	apps   []appPacket
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

type appPacket struct {
	AppName string `json:"app_name"`
	AppPath string `json:"app_path"`
}

func (s *Service) AddClient(clientID string, ws *cws.Client) *Client {
	client := NewServiceClient(clientID, ws)
	s.clients[clientID] = client
	return client
}

func (s *Service) RemoveClient(clientID string) {
	client := s.clients[clientID]
	delete(s.clients, clientID)
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
	delete(s.hosts, hostID)
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

//func addForwardingRoute(sender *cws.Client, senderID string, receiver *cws.Client, receiverID string, messages []string, s *Server, is_sender_browser bool) {
func addForwardingRoute(c *Client, h *Host, messages []string, s *Server, is_sender_browser bool) {
	for _, message := range messages {
		//sender.Receive(
		//	message,
		//	func(req cws.WSPacket) cws.WSPacket {
		//
		//		if (is_sender_browser) {
		//			if (s.capp.clientChosenApp[senderID].hostID != receiverID) {
		//				return cws.EmptyPacket
		//			}
		//		} else {
		//			if (s.capp.clientChosenApp[receiverID].hostID != senderID) {
		//				return cws.EmptyPacket
		//			}
		//		}
		//
		//		resp := receiver.SyncSend(req)
		//		return resp
		//	},
		//)
		var sender_ws *cws.Client
		var rec_ws *cws.Client
		if is_sender_browser {
			sender_ws = c.ws
			rec_ws = h.ws
		} else {
			sender_ws = h.ws
			rec_ws = c.ws
		}

		sender_ws.Receive(
			message,
			func(req cws.WSPacket) cws.WSPacket {
				resp := rec_ws.SyncSend(req)
				return resp
			})
	}
}

func GetAllHosts(s *Server) string {
	type MinimalHostMeta struct {
		HostID   string   `json:"host_id"`
		AppPaths []string `json:"app_paths"`
	}

	//type HostsMeta struct {
	//	hosts []MinimalHostMeta	`json:"hosts"`
	//}

	type HostsMeta []MinimalHostMeta
	var hosts = HostsMeta{}

	//var hosts []MinimalHostMeta
	for _, h := range s.capp.hosts {
		var paths = []string{}
		// Temp
		for _, app := range h.apps {
			paths = append(paths, app.AppPath)
		}
		hosts = append(hosts, MinimalHostMeta{h.hostID, paths})
	}

	//hostsData := HostsMeta{
	//	hosts: hosts,
	//}

	hostsJsonData, err := json.Marshal(hosts)
	if err != nil {
		return ""
	}

	hostJsonStr := string(hostsJsonData)
	return hostJsonStr
}

func (h *Host) HostRoute(s *Server) {
	h.ws.Receive(
		"registerApps",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			var registerAppsData struct {
				Apps []appPacket `json:"apps"`
			}
			log.Println("Get app registrations from host")
			err := json.Unmarshal([]byte(req.Data), &registerAppsData)
			if err != nil {
				log.Println("Error: Cannot decode json:", err)
				return cws.EmptyPacket
			}

			h.apps = registerAppsData.Apps

			hostsJsonStr := GetAllHosts(s)
			log.Println("hostsJsonStr " + hostsJsonStr)

			for _, client := range s.capp.clients {
				client.ws.Send(cws.WSPacket{
					Type: "hostsUpdated",
					Data: hostsJsonStr,
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
				HostID   string `json:"host_id"`
				AppParam string `json:"app"`
			}
			err := json.Unmarshal([]byte(req.Data), &hostAppsData)
			if err != nil {
				log.Println("Error: Cannot decode json:", err)
				return cws.EmptyPacket
			}
			chosen_host_app := ChosenHostApp{
				hostID:  hostAppsData.HostID,
				appPath: hostAppsData.AppParam,
			}
			s.capp.clientChosenApp[c.clientID] = chosen_host_app

			var h = s.capp.hosts[hostAppsData.HostID]

			if h != nil {
				addForwardingRoute(c, h, []string{"initwebrtc", "answer", "candidate"}, s, true)
				addForwardingRoute(c, h, []string{"init", "INIT", "candidate", "offer"}, s, false)

				// Send this request back to host. Host will start ffmpeg,
				// and send the request back to browser to start initwebrtc
				h.ws.Send(cws.WSPacket{
					Type: "init",
					Data: hostAppsData.AppParam,
				}, nil)
			}

			return cws.EmptyPacket
		},
	)
}

// NewCloudService returns a Cloud Service
func NewCloudService(cfg config.Config) *Service {
	s := &Service{
		clients:         map[string]*Client{},
		hosts:           map[string]*Host{},
		ccApp:           NewCloudAppClient(cfg),
		config:          cfg,
		clientChosenApp: map[string]ChosenHostApp{},
	}

	return s
}
