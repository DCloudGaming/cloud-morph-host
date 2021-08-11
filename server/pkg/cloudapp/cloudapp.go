// Package cloudapp is an individual cloud application
package cloudapp

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
)

type ccImpl struct {
	isReady       bool
	payloadType  uint8
}

// Packet represents a packet in cloudapp
type Packet struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// NewCloudAppClient returns new cloudapp client
func NewCloudAppClient(cfg config.Config) *ccImpl {
	c := &ccImpl{}

	log.Println("Launching application")
	c.launchApp(cfg);
	log.Println("Launched application")

	return c
}

func runApp(params []string) {
	log.Println("params: ", params)

	// Launch application using exec
	var cmd *exec.Cmd
	params = append([]string{"/C", "run-app.bat"}, params...)
	cmd = exec.Command("cmd", params...)

	cmd.Env = os.Environ()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	go func() {
		buf := bufio.NewReader(stdout) // Notice that this is not in a loop
		for {
			line, _, _ := buf.ReadLine()
			if string(line) == "" {
				continue
			}
			log.Println(string(line))
		}
	}()
	log.Println("execed run-client.sh")
	cmd.Wait()
}

// done to forcefully stop all processes
func (c *ccImpl) launchApp(cfg config.Config) chan struct{} {
	params := []string{cfg.Path, cfg.AppFile, cfg.WindowTitle}
	if cfg.HWKey {
		params = append(params, "game")
	} else {
		params = append(params, "")
	}
	params = append(params, []string{strconv.Itoa(cfg.ScreenWidth), strconv.Itoa(cfg.ScreenHeight)}...)

	runApp(params)

	done := make(chan struct{})

	return done
}
