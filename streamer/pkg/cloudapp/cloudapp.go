// Package cloudapp is an individual cloud application
package cloudapp

import (
	"bufio"
	"container/ring"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"

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
	// SendInput(Packet) TODO: Implement Input
	Handle()
	// TODO: this ssrc don't need to be exposed as interface
	GetSSRC() uint32
}

type ccImpl struct {
	isReady       bool
	videoListener *net.UDPConn
	audioListener *net.UDPConn
	videoStream   chan *rtp.Packet
	audioStream   chan *rtp.Packet
	inputEvents   chan Packet
	// gameConn      *net.TCPConn // talk with game
	screenWidth  float32
	screenHeight float32
	ssrc         uint32
	payloadType  uint8
	cfg config.Config
}

// Packet represents a packet in cloudapp
type Packet struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

const startVideoRTPPort = 5005
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
		cfg: cfg,
		//inputEvents: inputEvents,
	}

	if appPath != "" {
		NewCloudAppClientStart(c, appPath)
	}
	return c
}

func NewCloudAppClientStart(c *ccImpl, appPath string) {
	// To use for communicate with syncinput
	// la, err := net.ResolveTCPAddr("tcp4", ":9090")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("listening input at port 9090")
	// ln, err := net.ListenTCP("tcp", la)
	// if err != nil {
	// 	panic(err)
	// }

	log.Println("Launching application")
	c.launchApp(curVideoRTPPort, curAudioRTPPort, c.cfg, appPath)
	log.Println("Launched host app")

	// Read video stream from encoded video stream produced by FFMPEG
	log.Println("Setup Video Listener")
	videoListener, listenerssrc := c.newLocalStreamListener(curVideoRTPPort)
	c.videoListener = videoListener
	c.ssrc = listenerssrc
	log.Println("Setup Audio Listener")
	// TODO: Read video stream from encoded video stream produced by FFMPEG
	// audioListener, audiolistenerssrc := c.newLocalStreamListener(curAudioRTPPort)
	// c.audioListener = audioListener
	// c.ssrc = audiolistenerssrc

	c.listenVideoStream()
	log.Println("Launched Video stream listener")
	// c.listenAudioStream()
	// log.Println("Launched Audio stream listener")
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

func runApp(params []string, appPath string) {
	log.Println("params: ", params)

	// Launch application using exec
	var cmd *exec.Cmd
	params = append([]string{"/C", "bat_files/" + appPath}, params...)
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
func (c *ccImpl) launchApp(curVideoRTPPort int, curAudioRTPPort int, cfg config.Config, appPath string) chan struct{} {
	params := []string{cfg.Path, cfg.AppFile, cfg.WindowTitle}
	if cfg.HWKey {
		params = append(params, "game")
	} else {
		params = append(params, "")
	}
	params = append(params, []string{strconv.Itoa(cfg.ScreenWidth), strconv.Itoa(cfg.ScreenHeight)}...)

	runApp(params, appPath)
	// update flag
	c.screenWidth = float32(cfg.ScreenWidth)
	c.screenHeight = float32(cfg.ScreenHeight)

	done := make(chan struct{})

	return done
}

func (c *ccImpl) Handle() {
	// for event := range c.inputEvents {
	// c.SendInput(event)
	// }
}

// newLocalStreamListener returns RTP listener: listener and (Synchronization source) SSRC of that listener
func (c *ccImpl) newLocalStreamListener(rtpPort int) (*net.UDPConn, uint32) {
	// Open a UDP Listener for RTP Packets on port 5004
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("localhost"), Port: rtpPort})
	if err != nil {
		panic(err)
	}

	// Listen for a single RTP Packet, we need this to determine the SSRC
	inboundRTPPacket := make([]byte, 4096) // UDP MTU
	n, _, err := listener.ReadFromUDP(inboundRTPPacket)
	if err != nil {
		panic(err)
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
			c.videoListener.Close()
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

// func (c *ccImpl) SendInput(packet Packet) {
// 	switch packet.Type {
// 	case eventKeyUp:
// 		c.simulateKey(packet.Data, 0)
// 	case eventKeyDown:
// 		c.simulateKey(packet.Data, 1)
// 	case eventMouseMove:
// 		c.simulateMouseEvent(packet.Data, 0)
// 	case eventMouseDown:
// 		c.simulateMouseEvent(packet.Data, 1)
// 	case eventMouseUp:
// 		c.simulateMouseEvent(packet.Data, 2)
// 	}
// }

// func (c *ccImpl) simulateKey(jsonPayload string, keyState byte) {
// 	if !c.isReady {
// 		return
// 	}

// 	log.Println("KeyDown event", jsonPayload)
// 	type keydownPayload struct {
// 		KeyCode int `json:keycode`
// 	}
// 	p := &keydownPayload{}
// 	json.Unmarshal([]byte(jsonPayload), &p)

// 	vmKeyMsg := fmt.Sprintf("K%d,%b|", p.KeyCode, keyState)
// 	b, err := c.wineConn.Write([]byte(vmKeyMsg))
// 	log.Printf("%+v\n", c.wineConn)
// 	log.Println("Sended key: ", b, err)
// }

// // simulateMouseEvent handles mouse down event and send it to Virtual Machine over TCP port
// func (c *ccImpl) simulateMouseEvent(jsonPayload string, mouseState int) {
// 	if !c.isReady {
// 		return
// 	}

// 	type mousePayload struct {
// 		IsLeft byte    `json:isLeft`
// 		X      float32 `json:x`
// 		Y      float32 `json:y`
// 		Width  float32 `json:width`
// 		Height float32 `json:height`
// 	}
// 	p := &mousePayload{}
// 	json.Unmarshal([]byte(jsonPayload), &p)
// 	p.X = p.X * c.screenWidth / p.Width
// 	p.Y = p.Y * c.screenHeight / p.Height

// 	// Mouse is in format of comma separated "12.4,52.3"
// 	vmMouseMsg := fmt.Sprintf("M%d,%d,%f,%f,%f,%f|", p.IsLeft, mouseState, p.X, p.Y, p.Width, p.Height)
// 	_, err := c.wineConn.Write([]byte(vmMouseMsg))
// 	if err != nil {
// 		fmt.Println("Err: ", err)
// 	}
// }
