package main

import (
	"bytes"

	"log"
	"os"
	"os/signal"
	"time"

	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/net/websocket"
)

const (
	UPADTE = iota
	UPLOAD_ID
	ALL_DEVICES
	HEART_BEAT
	TAKE_PHOTO
)

type Command struct {
	CommandCode    int
	DeviceID       string
	CommandMessage string
}

var origin = "http://192.168.1.103/"
var url = "ws://192.168.1.103:23456/clientMainServer"

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	done := make(chan struct{})

	go func() {
		defer ws.Close()
		defer close(done)
		for {
			var command Command
			// Receive receives a text message serialized T as JSON.
			err := websocket.JSON.Receive(ws, &command)
			if err != nil {
				fmt.Println("clientMain:" + err.Error())
				return
			}

			switch command.CommandCode {
			case TAKE_PHOTO:

			}

		}
	}()

	LogIn(ws)

	ticker := time.NewTicker(time.Second * 100)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			command := Command{HEART_BEAT, "123456", t.String()}
			err := websocket.JSON.Send(ws, command)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := ws.WriteClose(1)
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func LogIn(ws *websocket.Conn) {
	getSerialNumber := "cat /proc/cpuinfo | grep Serial | awk ' {print $3}'"
	cmd := exec.Command("/bin/sh", "-c", getSerialNumber)
	var out bytes.Buffer //缓冲字节

	cmd.Stdout = &out //标准输出
	err := cmd.Run()  //运行指令 ，做判断
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String()) //输出执行结果
	serialNumber := out.String()

	serialNumber = strings.Replace(serialNumber, "\n", "", -1)
	command := Command{UPLOAD_ID, serialNumber, "this is device id"}
	websocket.JSON.Send(ws, command)
}
