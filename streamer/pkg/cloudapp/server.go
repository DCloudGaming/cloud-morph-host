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
	"net/http"
	"net/url"
)

type initData struct {
	CurAppID string `json:"cur_app_id"`
}

//const addr string = ":8081"

var signallingServerAddr = flag.String("addr", "localhost:8080", "http service address")

type Server struct {
	appID      string
	httpServer *http.Server
	wsClients  map[string]*cws.Client
	capp       *Service
}

type StreamerHttp struct {
	server *Server
}

//func (params *StreamerHttp) registerAppApi(w http.ResponseWriter, req *http.Request) {
//	fmt.Println("Receive Register App Requests")
//}

func NewServer(cfg config.Config) *Server {
	return NewServerWithHTTPServerMux(cfg)
}

func NewServerWithHTTPServerMux(cfg config.Config) *Server {
	//r := mux.NewRouter()
	//svmux := &http.ServeMux{}
	//svmux.Handle("/", r)

	//httpServer := &http.Server{
	//	Addr: addr,
	//	ReadTimeout: 5 * time.Second,
	//	WriteTimeout: 5 * time.Second,
	//	IdleTimeout: 120 * time.Second,
	//	Handler: svmux,
	//}
	server := &Server{
		capp: NewCloudService(cfg),
		//httpServer: httpServer,
	}

	//params := &StreamerHttp{server: server}
	//r.HandleFunc("/registerApp", params.registerAppApi)

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

func sendRegisterApp(s *Server) {
	// Send registrationApp Metadata to server
	for _, serviceClient := range s.capp.clients {

		type registerData struct {
			AppPaths []string `json:"app_paths"`
		}

		data := registerData{
			// TODO: User interact with GUI => Create those bat files. Then send
			// selection of which bat files allowed to use.
			AppPaths: []string{"run-notepad.bat", "run-chrome.bat"},
		}
		registerJsonData, err := json.Marshal(data)
		if err != nil {
			return
		}

		serviceClient.ws.Send(cws.WSPacket{
			Type: "registerApps",
			Data: string(registerJsonData),
			}, nil)
	}
}

func (s *Server) NotifySignallingServer() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *signallingServerAddr, Path: "/host"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	// Create websocket Client
	wsClient := cws.NewClient(c)
	clientID := wsClient.GetID()

	// Add new client game session to Cloud App service
	serviceClient := s.capp.AddClient(clientID, wsClient)
	s.initClientData(wsClient)

	serviceClient.Route(s.capp.GetSSRC(), s)
	log.Println("Initialized ServiceClient")

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

	// TODO: This function is invoked in registerAppApi, receive requests from GUI.
	sendRegisterApp(s)
}

func (o *Server) ListenAndServe() error {
	log.Println("Host http is running at", o.httpServer.Addr)
	return o.httpServer.ListenAndServe()
}

func (o *Server) Shutdown() {
}
