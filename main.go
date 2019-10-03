package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
	"golang.org/x/net/websocket"
	"time"
)

type Info struct {
	Title   string
	Message string
	Version string
}
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewInfo() *Info {
	result := &Info{
		Title:   "",
		Message: "智能语音客服助手",
		Version: "0.1.0",
	}
	return result
}
func basic() Info {
	info := Info{Message: "", Version: "0.0.1"}
	return info
}
func (i *Info) getMsg() string {
	return i.Message
}
func (i *Info) getVersion() string {
	return i.Version
}
func (i *Info) WailsInit(runtime *wails.Runtime) error {
	var wsConn *websocket.Conn
	ctx, cancel := context.WithCancel(context.Background())
	runtime.Events.On("start", func(optionalData ...interface{}) {
		fmt.Printf("%v start rec\n", time.Now().UnixNano())
		runtime.Events.Emit("NEWS", "录制开始", time.Now().Format("2006-01-02 15:04:05"))
		wsConn = checkAndLinkServer()
		go soundBiz(ctx, wsConn)
	})
	runtime.Events.On("end", func(optionalData ...interface{}) {
		fmt.Printf("%v end rec\n", time.Now().UnixNano())

		cancel()
		runtime.Events.Emit("NEWS", "录制结束", time.Now().Format("2006-01-02 15:04:05"))
	})
	return nil
}
func Login(u, p string) string {
	r := Result{}
	if u == "admin" && p == "admin" {
		r.Code = 200
		r.Msg = "登录成功"
		bs, _ := json.Marshal(&r)
		return string(bs)
	} else {
		r.Code = 400
		r.Msg = "登录失败"
		bs, _ := json.Marshal(&r)
		return string(bs)
	}
}

func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "智能语音客服助手",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(basic)
	app.Bind(Login)
	app.Bind(NewInfo())
	_ = app.Run()
}
