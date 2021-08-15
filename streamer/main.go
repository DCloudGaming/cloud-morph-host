package main

import (
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/DCloudGaming/cloud-morph-host/pkg/cloudapp"
)

const configFilePath = "./config.yaml"

func main() {
	cfg, err := config.ReadConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	server := cloudapp.NewServer(cfg)
	server.NotifySignallingServer()
	//server.Handle()
	//server.ListenAndServe()

	//go func() {
	//	err := server.ListenAndServe()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	select {
	case <-stop:
		log.Println("Received SIGTERM, Quiting")
		server.Shutdown()
	}
}

