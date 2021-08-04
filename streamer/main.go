package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/DCloudGaming/cloud-morph-host/pkg/cloudapp"
)

const configFilePath = "./config.yaml"

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	server := cloudapp.NewServer()
	server.Handle()

	server.NotifySignallingServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	select {
	case <-stop:
		log.Println("Received SIGTERM, Quiting")
		server.Shutdown()
	}
}

