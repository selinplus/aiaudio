package main

import (
	"context"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"golang.org/x/net/websocket"
)

var HOST = "129.211.114.37:4445"

var END_TAG = "{\"end\": true}"

func checkAndLinkServer() *websocket.Conn {
	conn, err := websocket.Dial("ws://"+HOST+"/ws", websocket.SupportedProtocolVersion, "http://"+HOST+"/ws")
	errCheck(err)
	fmt.Println("connect to server success")
	return conn
}

func soundBiz(ctx context.Context, wsConn *websocket.Conn) {
	var frameChan = make(chan []int16)
	sliceSize := 1280 * 10
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	framesPerBuffer := make([]int16, sliceSize)

	// init PortAudio
	err := portaudio.Initialize()
	errCheck(err)
	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
	errCheck(err)
	go sendToServer(frameChan, wsConn)
	go receiveFromServer(wsConn)
	for {
		select {
		case <-ctx.Done():
			_ = websocket.Message.Send(wsConn, END_TAG)
			_ = wsConn.Close()
			_ = stream.Close()
			_ = portaudio.Terminate()
			return
		default:
			errCheck(stream.Read())
			frameChan <- framesPerBuffer
		}
	}
}
func sendToServer(frameChan chan []int16, wsConn *websocket.Conn) {
	for {
		select {
		case fb := <-frameChan:
			b := make([]uint8, 0)
			for _, i16 := range fb {
				var h, l = uint8(i16 >> 8), uint8(i16 & 0xff)
				b = append(b, l)
				b = append(b, h)
			}
			err := websocket.Message.Send(wsConn, b)
			if string(b) == END_TAG {
				break
			}
			errCheck(err)
		}
	}
}
func receiveFromServer(wsConn *websocket.Conn) {
	var msg []byte
	for {
		err := websocket.Message.Receive(wsConn, &msg)
		errCheck(err)
		fmt.Println(string(msg))
		if string(msg) == "END" {
			break
		}
	}
}
