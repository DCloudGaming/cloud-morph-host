// Widget server to serve a standalone cloudmorph instance
package cloudapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/cws"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type initData struct {
	CurAppID string `json:"cur_app_id"`
}

const embedPage string = "web/index.html"
const addr string = ":8080"

type Server struct {
	appID      string
	httpServer *http.Server
	wsClients  map[string]*cws.Client
	capp       *Service
}

func NewServer(cfg config.Config) *Server {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web"))))

	svmux := &http.ServeMux{}
	svmux.Handle("/", r)

	return NewServerWithHTTPServerMux(cfg, r, svmux)
}

func NewServerWithHTTPServerMux(cfg config.Config, r *mux.Router, svmux *http.ServeMux) *Server {
	server := &Server{}

	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := template.ParseFiles(embedPage)
			if err != nil {
				log.Fatal(err)
			}

			tmpl.Execute(w, nil)
		},
	)
	// Websocket
	r.HandleFunc("/client", server.Client)
	r.HandleFunc("/host", server.Host)
	httpServer := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      svmux,
	}
	server.capp = NewCloudService(cfg)
	server.httpServer = httpServer

	return server
}

func (s *Server) Host(w http.ResponseWriter, r *http.Request) {
	log.Println("A host is connecting...")
	defer func() {
		if r := recover(); r != nil {
			log.Println("Warn: Something wrong. Recovered in ", r)
		}
	}()

	// upgrader to upgrade http connection to websocket connection
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// Check origin of upgrader
		// TODO: can we be stricter?
		return true
	}
	// be aware of ReadBufferSize, WriteBufferSize (default 4096)
	// https://pkg.go.dev/github.com/gorilla/websocket?tab=doc#Upgrader
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Coordinator: [!] WS upgrade:", err)
		return
	}

	// Create websocket Client
	wsClient := cws.NewClient(c)
	clientID := wsClient.GetID()
	// Add new client game session to Cloud App service
	s.capp.AddHost(clientID, wsClient)

	// TODO: add mapping host-client here
	for _, serviceClient := range s.capp.clients {
		addForwardingRoute(serviceClient.ws, wsClient, []string{"initwebrtc", "answer", "candidate"})
		addForwardingRoute(wsClient, serviceClient.ws, []string{"init", "INIT", "candidate", "offer"})
		break
	}

	log.Println("Initialized ServiceHost")

	go wsClient.Heartbeat()

	s.initClientData(wsClient)
	go func(browserClient *cws.Client) {
		browserClient.Listen()
		log.Println("Closing connection")
		browserClient.Close()
		s.capp.RemoveClient(clientID)
		log.Println("Closed connection")
	}(wsClient)
}

func (s *Server) Client(w http.ResponseWriter, r *http.Request) {
	log.Println("A user is connecting...")
	defer func() {
		if r := recover(); r != nil {
			log.Println("Warn: Something wrong. Recovered in ", r)
		}
	}()

	// upgrader to upgrade http connection to websocket connection
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// Check origin of upgrader
		// TODO: can we be stricter?
		return true
	}
	// be aware of ReadBufferSize, WriteBufferSize (default 4096)
	// https://pkg.go.dev/github.com/gorilla/websocket?tab=doc#Upgrader
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Coordinator: [!] WS upgrade:", err)
		return
	}

	// Create websocket Client
	wsClient := cws.NewClient(c)
	clientID := wsClient.GetID()
	// Add new client game session to Cloud App service
	s.capp.AddClient(clientID, wsClient)

	for _, serviceHost := range s.capp.hosts {
		addForwardingRoute(wsClient, serviceHost.ws, []string{"initwebrtc", "answer", "candidate"})
		addForwardingRoute(serviceHost.ws, wsClient, []string{"init", "INIT", "candidate", "offer"})
		break
	}
	log.Println("Initialized ServiceClient")

	s.initClientData(wsClient)
	go func(browserClient *cws.Client) {
		browserClient.Listen()
		log.Println("Closing connection")
		browserClient.Close()
		s.capp.RemoveClient(clientID)
		log.Println("Closed connection")
	}(wsClient)
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

func (o *Server) ListenAndServe() error {
	log.Println("Server is running at", addr)
	return o.httpServer.ListenAndServe()
}

func (o *Server) Shutdown() {
}
