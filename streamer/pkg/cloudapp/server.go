// Widget server to serve a standalone cloudmorph instance
package cloudapp

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/cws"

	//"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	//"time"
)

type initData struct {
	CurAppID string `json:"cur_app_id"`
}

const addr string = ":8082"

type Server struct {
	appID            string
	httpServer       *http.Server
	wsClients        map[string]*cws.Client
	capp             *Service
	signalServerAddr string
	token            string
}

type StreamerHttp struct {
	server *Server
}

type appPacket struct {
	AppName string `json:"app_name"`
	AppPath string `json:"app_path"`
}

func (params *StreamerHttp) registerAppApi(w http.ResponseWriter, req *http.Request) {
	fmt.Println("received register package")

	var appBody []appPacket
	json.NewDecoder(req.Body).Decode(&appBody)

	// TODO: Create .bat file contents here
	sendRegisterApp(params.server, appBody)
}

func sendRegisterApp(s *Server, appBody []appPacket) {
	// Send registrationApp Metadata to server
	for _, serviceClient := range s.capp.clients {

		type registerData struct {
			Apps []appPacket `json:"apps"`
		}

		data := registerData{
			// selection of which bat files allowed to use.
			//AppPaths: []string{"run-notepad.bat", "run-chrome.bat"},
			Apps: appBody,
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

type tokenPacket struct {
	Token string `json:"token"`
}

func (params *StreamerHttp) updateTokenApi(w http.ResponseWriter, req *http.Request) {
	fmt.Println("received register package")

	var tokenBody tokenPacket
	json.NewDecoder(req.Body).Decode(&tokenBody)

	updateToken(params.server, tokenBody)
}

func updateToken(s *Server, tokenBody tokenPacket) {
	for _, serviceClient := range s.capp.clients {

		updateTokenData, err := json.Marshal(tokenBody)
		if err != nil {
			return
		}

		serviceClient.ws.Send(cws.WSPacket{
			Type: "updateToken",
			Data: string(updateTokenData),
		}, nil)
	}
}

func NewServer(cfg config.Config) *Server {
	return NewServerWithHTTPServerMux(cfg)
}

func NewServerWithHTTPServerMux(cfg config.Config) *Server {
	//r := mux.NewRouter()
	//svmux := &http.ServeMux{}
	//svmux.Handle("/", r)
	//
	//httpServer := &http.Server{
	//	Addr: addr,
	//	ReadTimeout: 5 * time.Second,
	//	WriteTimeout: 5 * time.Second,
	//	IdleTimeout: 120 * time.Second,
	//	Handler: svmux,
	//}
	server := &Server{
		capp:             NewCloudService(cfg),
		signalServerAddr: "localhost:8080", // Default, will be overriden by flag
		//httpServer: httpServer,
	}
	flag.StringVar(&server.signalServerAddr, "addr", os.Getenv("HOST"), "http service address")

	params := &StreamerHttp{server: server}
	http.HandleFunc("/registerApp", params.registerAppApi)
	http.HandleFunc("/updateToken", params.updateTokenApi)
	//.Host("http://localhost:8081").Methods("GET").Schemes("http")

	go http.ListenAndServe(":8082", nil)
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

	u := url.URL{Scheme: "ws", Host: s.signalServerAddr, Path: "/host"}
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
	//sendRegisterApp(s)
}

//func (o *Server) ListenAndServe() error {
//	log.Println("Host http is running at", o.httpServer.Addr)
//	//err := o.httpServer.ListenAndServe()
//}

func (o *Server) Shutdown() {
}
