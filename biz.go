package main

import (
	"context"
	"github.com/gordonklaus/portaudio"
	"golang.org/x/net/websocket"
)

var HOST = "120.27.22.145"
var frameChan = make(chan []int16)

func sound(ctx context.Context) {

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
	go sendToServer()
	go receiveFromServer()
	for {
		select {
		case <-ctx.Done():
			_ = stream.Close()
			_ = portaudio.Terminate()
			return
		default:
			errCheck(stream.Read())
			frameChan <- framesPerBuffer
		}
	}
}
func sendToServer() {
	for {
		select {
		case fb := <-frameChan:
			err := websocket.Message.Send(wsConn, fb)
			errCheck(err)
		}
	}
}
func receiveFromServer() {
	var res map[string]string
	for {
		err := websocket.Message.Receive(wsConn, res)
		errCheck(err)
	}
}
