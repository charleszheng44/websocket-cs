package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	msg "github.com/charleszheng44/websocket-server/pkg/message"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

var (
	port        int
	uploadDir   string
	downloadDir string
)

func init() {
	flag.IntVar(&port, "port", 8898, "port client will dail to")
	flag.StringVar(&uploadDir, "uploaddir", "",
		"files in this directory will be uploaded")
	flag.StringVar(&downloadDir, "downloaddir", "",
		"recieved files will be stored in this directory")
	flag.Parse()
	if uploadDir == "" {
		logrus.Error("[CLIENT] please specify the upload directory using --uploaddir")
		os.Exit(1)
	}
	if downloadDir == "" {
		logrus.Error("[CLIENT] please specify the download directory using --downloaddir")
		os.Exit(1)
	}
}

// Client.
func main() {
	url := fmt.Sprintf("ws://localhost:%d/", port)
	origin := fmt.Sprintf("http://localhost:%d", port)
	ws, err := websocket.Dial(url, "", origin)

	// receive message
	go func() {
		for {
			var m msg.Message
			err = websocket.JSON.Receive(ws, &m)
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			logrus.Printf("[CLIENT] received message: %d\n", m.Id)
			log.Printf("[CLIENT] received messag: %d\n", m.Id)
			if err = m.GenFile(downloadDir); err != nil {
				logrus.Errorf("[CLIENT] fail to generate files: %v", err)
			}
		}
	}()

	// send messages
	files, err := ioutil.ReadDir(uploadDir)
	for i, f := range files {
		m, err := msg.GenMessage(i, f, uploadDir)
		if err != nil {
			logrus.Errorf("[CLIENT] fail to generate message: %s", err)
			os.Exit(1)
		}
		err = websocket.JSON.Send(ws, m)
		if err != nil {
			logrus.Errorf("[CLIENT] fail to send message: %s", err)
			os.Exit(1)
		}

		time.Sleep(2 * time.Second)
	}

	logrus.Println("[CLIENT] server finished request...")
}
