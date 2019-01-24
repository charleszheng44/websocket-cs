package util

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var connectionUpgradeRegex = regexp.MustCompile("(^|.*,\\s*)upgrade($|\\s*,)")

func IsWebsocketRequest(req *http.Request) bool {
	return connectionUpgradeRegex.MatchString(
		strings.ToLower(req.Header.Get("Connection"))) &&
		strings.ToLower(req.Header.Get("Upgrade")) == "websocket"
}

// CompressFile compresses the given file and write
// the output to the buffer
func CompressFile(path string, buf *bytes.Buffer) error {
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	byteArr, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	gw.Write(byteArr)
	return nil
}

// DecompressFile decompresses the file content hold
// by buffer to specified file
func DecompressFile(inBuf *bytes.Buffer, path string) error {
	gr, err := gzip.NewReader(inBuf)
	var outBuf bytes.Buffer
	if _, err = io.Copy(&outBuf, gr); err != nil {
		return err
	}
	gr.Close()

	outFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	defer outFile.Close()
	if err != nil {
		return err
	}
	if _, err = outFile.Write(outBuf.Bytes()); err != nil {
		return err
	}
	return nil
}
