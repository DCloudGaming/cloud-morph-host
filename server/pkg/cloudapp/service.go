package cloudapp

import (
	"encoding/json"
	"fmt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/jwt"
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
	walletAddress string
	ws       *cws.Client
	// cancel to trigger cleaning up when client is disconnected
	cancel chan struct{}
	// done to notify if the client is done clean up
	done chan struct{}
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

func addForwardingRoute(c *Client, h *Host, messages []string, s *Server, is_sender_browser bool) {
	for _, message := range messages {
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

	h.ws.Receive(
		"updateToken",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			var updateTokenData struct {
				Token string `json:"token"`
			}
			err := json.Unmarshal([]byte(req.Data), &updateTokenData)
			if err != nil {
				log.Println("Error: Cannot decode json:", err)
				return cws.EmptyPacket
			}

			user, _ := jwt.DecodeUser(updateTokenData.Token)
			h.walletAddress = user.WalletAddress
			fmt.Println("Set wallet address for host ws " + h.walletAddress)
			return cws.EmptyPacket
		},
	)
}

func (c *Client) ClientRoute(s *Server) {
	c.ws.Receive("startSession",
		func(req cws.WSPacket) (resp cws.WSPacket) {
			var startSessionData struct {
				AppName string `json:"app_name"`
				HostWalletAddress string `json:"host_wallet_address"`
			}
			err := json.Unmarshal([]byte(req.Data), &startSessionData)
			if err != nil {
				log.Println("Error: Cannot decode json:", err)
				return cws.EmptyPacket
			}

			for _, h := range s.capp.hosts {
				if (h.walletAddress == startSessionData.HostWalletAddress || s.shared_env.Mode() == "DEBUG") {
					addForwardingRoute(c, h, []string{"initwebrtc", "answer", "candidate"}, s, true)
					addForwardingRoute(c, h, []string{"init", "INIT", "candidate", "offer"}, s, false)

					var appPath string
					if (s.shared_env.Mode() == "DEBUG") {
						appPath = s.shared_env.DefaultAppPath()
					} else {
						registeredApp, _ := s.shared_env.AppRepo().GetAppByName(startSessionData.AppName, h.walletAddress)
						appPath = registeredApp.AppPath
					}

					h.ws.Send(cws.WSPacket{
						Type: "init",
						Data: appPath,
					}, nil)
					break
				}
			}

			return cws.EmptyPacket
		},
	)

	// TODO: Add forwarding mapping to an api handler
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
