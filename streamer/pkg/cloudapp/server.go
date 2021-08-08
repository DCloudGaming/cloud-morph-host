// Widget server to serve a standalone cloudmorph instance
package cloudapp

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/cws"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

type initData struct {
	CurAppID string `json:"cur_app_id"`
}

var signallingServerAddr = flag.String("addr", "localhost:8080", "http service address")

type Server struct {
	appID      string
	wsClients  map[string]*cws.Client
	capp       *Service
}

func NewServer(cfg config.Config) *Server {
	return NewServerWithHTTPServerMux(cfg)
}

func NewServerWithHTTPServerMux(cfg config.Config) *Server {
	server := &Server{
		capp: NewCloudService(cfg),
	}
	return server
}

func (o *Server) Handle() {
	// Spawn CloudGaming Handle
	go o.capp.Handle()
}

func (s *Server) initClientData(client *cws.Client) {
	data := initData{
		CurAppID: s.appID,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Println("Send Client INIT")
	client.Send(cws.WSPacket{
		Type: "INIT",
		Data: string(jsonData),
	}, nil)
}

func (s *Server) NotifySignallingServer() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *signallingServerAddr, Path: "/host"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Create websocket Client
	wsClient := cws.NewClient(c)
	clientID := wsClient.GetID()

	// Add new client game session to Cloud App service
	serviceClient := s.capp.AddClient(clientID, wsClient)
	serviceClient.Route(s.capp.GetSSRC())
	log.Println("Initialized ServiceClient")

	s.initClientData(wsClient)

	go func(browserClient *cws.Client) {
		browserClient.Listen()
		log.Println("Closing connection")
		browserClient.Close()
		s.capp.RemoveClient(clientID)
		log.Println("Closed connection")
	}(wsClient)

	if err != nil {
		log.Println("Coordinator: [!] WS upgrade:", err)
		return
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
}

func (o *Server) Shutdown() {
}
