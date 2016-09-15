package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"golang.org/x/net/websocket"
)

const (
	defaultPort = "8753"
)

var (
	ngmr io.ReadCloser
	ngmw io.WriteCloser
)

// BridgeServer the data received on the WebSocket.
func BridgeServer(wsc *websocket.Conn) {
	go func() {
		_, err := io.Copy(wsc, ngmr)
		if err != nil {
			log.Println(err)
		}
	}()
	_, err := io.Copy(ngmw, wsc)
	if err != nil {
		log.Println(err)
	}
}

// This example demonstrates a trivial echo server.
func main() {
	var err error

	// connect to Nagome
	cmd := exec.Command("nagome")
	ngmw, err = cmd.StdinPipe()
	if err != nil {
		log.Println(err)
		return
	}
	defer ngmw.Close()
	ngmr, err = cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	defer ngmr.Close()
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("http://localhost:" + defaultPort + "/app")

	// serve
	http.Handle("/ws", websocket.Handler(BridgeServer))
	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./app"))))
	err = http.ListenAndServe(":"+defaultPort, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
