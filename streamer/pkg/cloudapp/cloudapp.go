// Package cloudapp is an individual cloud application
package cloudapp

import (
	"bufio"
	"container/ring"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/DCloudGaming/cloud-morph-host/pkg/common/config"
	"github.com/DCloudGaming/cloud-morph-host/pkg/common/cws"
	"github.com/pion/rtp"
)

type InputEvent struct {
	inputType    bool
	inputPayload []byte
}

type CloudAppClient interface {
	VideoStream() chan *rtp.Packet
	// AudioStream() chan *rtp.Packet
	SendInput(Packet)
	Handle()
	// TODO: this ssrc don't need to be exposed as interface
	GetSSRC() uint32
}

type osTypeEnum int

const (
	Linux osTypeEnum = iota
	Mac
	Windows
)

type ccImpl struct {
	isReady       bool
	videoListener *net.UDPConn
	audioListener *net.UDPConn
	videoStream   chan *rtp.Packet
	audioStream   chan *rtp.Packet
	inputEvents   chan Packet
	// connection with syncinput script
	syncInputConn *net.TCPConn
	osType        osTypeEnum
	// gameConn      *net.TCPConn // talk with game
	screenWidth  float32
	screenHeight float32
	ssrc         uint32
	payloadType  uint8
	cfg          config.Config
}

// Packet represents a packet in cloudapp
type Packet struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

const startVideoRTPPort = 5006
const startAudioRTPPort = 4004
const eventKeyDown = "KEYDOWN"
const eventKeyUp = "KEYUP"
const eventMouseMove = "MOUSEMOVE"
const eventMouseDown = "MOUSEDOWN"
const eventMouseUp = "MOUSEUP"

var curVideoRTPPort = startVideoRTPPort
var curAudioRTPPort = startAudioRTPPort

// NewCloudAppClient returns new cloudapp client
func NewCloudAppClient(cfg config.Config, inputEvents chan Packet, appPath string) *ccImpl {
	c := &ccImpl{
		videoStream: make(chan *rtp.Packet, 1),
		audioStream: make(chan *rtp.Packet, 1),
		cfg:         cfg,
		inputEvents: inputEvents,
	}

	switch runtime.GOOS {
	case "windows":
		c.osType = Windows
	default:
		c.osType = Linux
	}

	// NewCloudAppClientStart(c, "/home/giongto/Desktop/code/cloud-morph-host/streamer/apps/Minesweeper.exe")
	if appPath != "" {
		NewCloudAppClientStart(c, appPath)
	}
	return c
}

func NewCloudAppClientStart(c *ccImpl, appPath string) {
	// To use for communicate with syncinput
	la, err := net.ResolveTCPAddr("tcp4", ":9090")
	if err != nil {
		panic(err)
	}
	log.Println("listening input at port 9090")
	ln, err := net.ListenTCP("tcp", la)
	if err != nil {
		panic(err)
	}

	log.Println("Launching application")
	c.launchApp(curVideoRTPPort, curAudioRTPPort, c.cfg, appPath)
	log.Println("Launched host app")

	// Read video stream from encoded video stream produced by FFMPEG
	log.Println("Setup Video Listener")
	videoListener, listenerssrc := c.newLocalStreamListener(curVideoRTPPort)
	c.videoListener = videoListener
	c.ssrc = listenerssrc
	if c.osType != Windows {
		// Don't spawn Audio in Windows
		log.Println("Setup Audio Listener")
		audioListener, audiolistenerssrc := c.newLocalStreamListener(curAudioRTPPort)
		c.audioListener = audioListener
		c.ssrc = audiolistenerssrc
	}
	log.Println("Setup Audio Listener")
	// TODO: Read video stream from encoded video stream produced by FFMPEG
	// audioListener, audiolistenerssrc := c.newLocalStreamListener(curAudioRTPPort)
	// c.audioListener = audioListener
	// c.ssrc = audiolistenerssrc

	c.listenVideoStream()
	log.Println("Launched Video stream listener")
	if c.osType != Windows {
		// Don't spawn Audio in Windows
		c.listenAudioStream()
		log.Println("Launched Audio stream listener")
	}
	// c.listenAudioStream()
	// log.Println("Launched Audio stream listener")

	// Maintain input stream from server to Virtual Machine over websocket
	go c.healthCheckVM()
	// NOTE: Why Websocket: because normal IPC cannot communicate cross OS.
	go func() {
		for {
			log.Println("Waiting syncinput to connect")
			// Polling Wine socket connection (input stream)
			conn, err := ln.AcceptTCP()
			log.Println("Accepted a TCP connection")
			if err != nil {
				log.Println("err: ", err)
			}
			conn.SetKeepAlive(true)
			conn.SetKeepAlivePeriod(10 * time.Second)
			c.syncInputConn = conn
			c.isReady = true
			log.Println("Launched IPC with VM")
		}
	}()

}

// convertWSPacket returns cloudapp packet from ws packet
func convertWSPacket(packet cws.WSPacket) Packet {
	return Packet{
		Type: packet.Type,
		Data: packet.Data,
	}
}

func (c *ccImpl) GetSSRC() uint32 {
	return c.ssrc
}

func (c *ccImpl) runApp(params []string) {

	// Launch application using exec
	var cmd *exec.Cmd
	//params = append([]string{"/C", "run-app.bat"}, params...)
	if c.osType == Windows {
		params = append([]string{"-ExecutionPolicy", "Bypass", "-F", "run-app.ps1"}, params...)
		log.Println("You are running on Windows", params)
		cmd = exec.Command("powershell", params...)
	} else {
		log.Println("You are running on Linux")
		cmd = exec.Command("./run-wine.sh", params...)
	}

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
	// cmd.Wait()
}

// done to forcefully stop all processes
func (c *ccImpl) launchApp(curVideoRTPPort int, curAudioRTPPort int, cfg config.Config, appPath string) chan struct{} {
	_, filename := filepath.Split(appPath)
	fmt.Println("Running ", appPath, " ", filename)
	params := []string{}
	if c.osType == Windows {
		params = append(params, []string{appPath, filename}...)
	} else {
		mountedDirName := "/winevm/apps"
		appTitle := "Minesweeper"
		params = append(params, []string{mountedDirName, filename, appTitle}...)
	}

	if cfg.HWKey {
		params = append(params, "game")
	} else {
		params = append(params, "")
	}
	params = append(params, []string{strconv.Itoa(cfg.ScreenWidth), strconv.Itoa(cfg.ScreenHeight), appPath}...)

	log.Println("params: ", params)
	c.runApp(params)
	// update flag
	c.screenWidth = float32(cfg.ScreenWidth)
	c.screenHeight = float32(cfg.ScreenHeight)

	done := make(chan struct{})

	return done
}

// healthCheckVM to maintain connection with Virtual Machine
func (c *ccImpl) healthCheckVM() {
	log.Println("Starting health check")
	for {
		if c.syncInputConn != nil {
			_, err := c.syncInputConn.Write([]byte{0})
			if err != nil {
				log.Println(err)
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func (c *ccImpl) Handle() {
	for event := range c.inputEvents {
		c.SendInput(event)
	}
}

// newLocalStreamListener returns RTP listener: listener and (Synchronization source) SSRC of that listener
func (c *ccImpl) newLocalStreamListener(rtpPort int) (*net.UDPConn, uint32) {
	fmt.Println("new Local Stream")
	// Open a UDP Listener for RTP Packets on port 5006
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("localhost"), Port: rtpPort})
	if err != nil {
		fmt.Errorf("%v", err)
		// panic(err)
		return nil, 0
	}

	// Listen for a single RTP Packet, we need this to determine the SSRC
	inboundRTPPacket := make([]byte, 4096) // UDP MTU
	n, _, err := listener.ReadFromUDP(inboundRTPPacket)
	if err != nil {
		fmt.Errorf("%v", err)
		// panic(err)
		return nil, 0
	}

	// Unmarshal the incoming packet
	packet := &rtp.Packet{}
	if err = packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
		panic(err)
	}

	return listener, packet.SSRC
}

func (c *ccImpl) VideoStream() chan *rtp.Packet {
	return c.videoStream
}

func (c *ccImpl) AudioStream() chan *rtp.Packet {
	return c.audioStream
}

// Listen to videostream, output to videoStream channel
func (c *ccImpl) listenAudioStream() {

	// Broadcast video stream
	go func() {
		defer func() {
			c.audioListener.Close()
			log.Println("Closing app VM")
		}()
		r := ring.New(120)

		n := r.Len()
		for i := 0; i < n; i++ {
			// r.Value = make([]byte, 4096)
			r.Value = make([]byte, 1500)
			r = r.Next()
		}

		// TODO: Create a precreated memory, only pop after finish processing
		// Read RTP packets forever and send them to the WebRTC Client
		for {
			inboundRTPPacket := r.Value.([]byte) // UDP MTU
			r = r.Next()
			n, _, err := c.audioListener.ReadFrom(inboundRTPPacket)
			if err != nil {
				log.Printf("error during read: %s", err)
				continue
			}

			// TODOs: Don't assign packet here
			packet := &rtp.Packet{}
			if err := packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
				log.Printf("error during unmarshalling a packet: %s", err)
				continue
			}

			c.audioStream <- packet
		}
	}()

}

// Listen to videostream, output to videoStream channel
func (c *ccImpl) listenVideoStream() {

	// Broadcast video stream
	go func() {
		defer func() {
			// c.videoListener.Close()
			log.Println("Closing app VM")
		}()
		r := ring.New(120)

		n := r.Len()
		for i := 0; i < n; i++ {
			r.Value = make([]byte, 1500)
			r = r.Next()
		}

		// Read RTP packets forever and send them to the WebRTC Client
		for {
			inboundRTPPacket := r.Value.([]byte) // UDP MTU
			r = r.Next()
			if c.videoListener == nil {
				continue
			}
			n, _, err := c.videoListener.ReadFrom(inboundRTPPacket)
			if err != nil {
				log.Printf("error during read: %s", err)
				continue
			}

			// TODOs: Don't assign packet here
			packet := &rtp.Packet{}
			if err := packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
				log.Printf("error during unmarshalling a packet: %s", err)
				continue
			}

			c.videoStream <- packet
		}
	}()

}

func (c *ccImpl) SendInput(packet Packet) {
	switch packet.Type {
	case eventKeyUp:
		c.simulateKey(packet.Data, 0)
	case eventKeyDown:
		c.simulateKey(packet.Data, 1)
	case eventMouseMove:
		c.simulateMouseEvent(packet.Data, 0)
	case eventMouseDown:
		c.simulateMouseEvent(packet.Data, 1)
	case eventMouseUp:
		c.simulateMouseEvent(packet.Data, 2)
	}
}

func (c *ccImpl) simulateKey(jsonPayload string, keyState byte) {
	if !c.isReady {
		return
	}

	log.Println("KeyDown event", jsonPayload)
	type keydownPayload struct {
		KeyCode int `json:keycode`
	}
	p := &keydownPayload{}
	json.Unmarshal([]byte(jsonPayload), &p)

	vmKeyMsg := fmt.Sprintf("K%d,%b|", p.KeyCode, keyState)
	b, err := c.syncInputConn.Write([]byte(vmKeyMsg))
	log.Printf("%+v\n", c.syncInputConn)
	log.Println("Sended key: ", b, err)
}

// simulateMouseEvent handles mouse down event and send it to Virtual Machine over TCP port
func (c *ccImpl) simulateMouseEvent(jsonPayload string, mouseState int) {
	if !c.isReady {
		return
	}

	type mousePayload struct {
		IsLeft byte    `json:isLeft`
		X      float32 `json:x`
		Y      float32 `json:y`
		Width  float32 `json:width`
		Height float32 `json:height`
	}
	p := &mousePayload{}
	json.Unmarshal([]byte(jsonPayload), &p)
	p.X = p.X * c.screenWidth / p.Width
	p.Y = p.Y * c.screenHeight / p.Height

	// Mouse is in format of comma separated "12.4,52.3"
	vmMouseMsg := fmt.Sprintf("M%d,%d,%f,%f,%f,%f|", p.IsLeft, mouseState, p.X, p.Y, p.Width, p.Height)
	_, err := c.syncInputConn.Write([]byte(vmMouseMsg))
	if err != nil {
		fmt.Println("Err: ", err)
	}
}
