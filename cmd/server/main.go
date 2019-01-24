package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	msg "github.com/charleszheng44/websocket-server/pkg/message"
	util "github.com/charleszheng44/websocket-server/pkg/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

var (
	port        int
	uploadDir   string
	downloadDir string
)

func init() {
	flag.IntVar(&port, "port", 8898, "port the server will listen at")
	flag.StringVar(&uploadDir, "uploaddir", "",
		"files in this directory will be uploaded")
	flag.StringVar(&downloadDir, "downloaddir", "",
		"recieved files will be stored in this directory")
	flag.Parse()
	if uploadDir == "" {
		logrus.Error("[SERVER] please specify the upload directory using --uploaddir")
		os.Exit(1)
	}
	if downloadDir == "" {
		logrus.Error("[SERVER] please specify the download directory using --downloaddir")
		os.Exit(1)
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	// Handle websockets if specified.
	if util.IsWebsocketRequest(r) {
		websocket.Handler(HandleWebSockets).ServeHTTP(w, r)
	} else {
		log.Fatal("Only support websocket protocol")
	}
	log.Println("Finished sending response...")
}

func HandleWebSockets(ws *websocket.Conn) {
	// receive message from client
	go func() {
		for {
			var m msg.Message
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			log.Printf("[SERVER] received messag: %d\n", m.Id)
			if err = m.GenFile(downloadDir); err != nil {
				logrus.Errorf("[SERVER] fail to generate files: %v", err)
			}
		}
	}()

	// loop the files in uploadDir
	files, err := ioutil.ReadDir(uploadDir)
	if err != nil {
		logrus.Errorf("[SERVER] fail to read files from %s: %s",
			uploadDir, err)
		os.Exit(1)
	}
	for i, f := range files {
		// generate messages for each file
		m, err := msg.GenMessage(i, f, uploadDir)
		if err != nil {
			logrus.Errorf("[SERVER] fail to gen message: %s", err)
			os.Exit(1)
		}
		err = websocket.JSON.Send(ws, m)
		if err != nil {
			logrus.Errorf("[SERVER] fail to send message: %s", err)
			os.Exit(1)
		}
		// force the server to sleep
		time.Sleep(2 * time.Second)
	}

	// TODO check how to use mutex.Wait
}

// Server.
func main() {
	http.HandleFunc("/", Handle)
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Serving at %s ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
