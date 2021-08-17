// Widget server to serve a standalone cloudmorph instance
package cloudapp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/cws"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type initData struct {
	CurAppID string `json:"cur_app_id"`
}

const embedPageIndex string = "web/index.html"
const embedPagePlay string = "web/play.html"
const embedPageRegister string = "web/register.html"
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
			tmpl, err := template.ParseFiles(embedPageIndex)
			if err != nil {
				log.Fatal(err)
			}

			tmpl.Execute(w, nil)
		},
	)
	r.HandleFunc("/register",
		func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := template.ParseFiles(embedPageRegister)
			if err != nil {
				log.Fatal(err)
			}

			tmpl.Execute(w, nil)
		},
	)
	r.HandleFunc("/play",
		func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := template.ParseFiles(embedPagePlay)
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

	// Create websocket client for Host
	wsHost := cws.NewClient(c)
	hostID := wsHost.GetID()
	// Add new client game session to Cloud App service
	serviceHostClient := s.capp.AddHost(hostID, wsHost)

	// TODO: add mapping host-client here. This one should be invoked when client requests particular one.
	//for _, serviceClient := range s.capp.clients {
	//	addForwardingRoute(serviceClient.ws, serviceClient.clientID, wsHost, hostID, []string{"initwebrtc", "answer", "candidate"}, s, true)
	//	addForwardingRoute(wsHost, hostID, serviceClient.ws, serviceClient.clientID, []string{"init", "INIT", "candidate", "offer"}, s, false)
	//	break
	//}

	log.Println("Initialized ServiceHost")

	go wsHost.Heartbeat()

	s.initClientData(wsHost)
	serviceHostClient.HostRoute(s)

	go func(hostClient *cws.Client) {
		hostClient.Listen()
		log.Println("Closing connection")
		hostClient.Close()
		s.capp.RemoveHost(hostID)
		log.Println("Closed connection")
	}(wsHost)
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
	serviceBrowserClient := s.capp.AddClient(clientID, wsClient)

	//for _, serviceHost := range s.capp.hosts {
	//	addForwardingRoute(wsClient, clientID, serviceHost.ws, serviceHost.hostID, []string{"initwebrtc", "answer", "candidate"}, s, true)
	//	addForwardingRoute(serviceHost.ws, serviceHost.hostID, wsClient, clientID, []string{"init", "INIT", "candidate", "offer"}, s, false)
	//	break
	//}
	log.Println("Initialized ServiceClient")

	s.initClientData(wsClient)
	serviceBrowserClient.ClientRoute(s)

	go func(browserClient *cws.Client) {
		browserClient.Listen()
		log.Println("Closing connection")
		browserClient.Close()
		s.capp.RemoveClient(clientID)
		log.Println("Closed connection")
	}(wsClient)

	wsClient.Send(cws.WSPacket{
		Type: "hostsUpdated",
		Data: GetAllHosts(s),
	}, nil)
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
