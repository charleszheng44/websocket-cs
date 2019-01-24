package message

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charleszheng44/websocket-server/pkg/util"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Id           int    `json:"id,omitempty"`
	Name         string `json:"Name,omitempty"`
	InputPath    string `json:"inputpath,omitempty"`
	InputContent []byte `json:"inputcontent,omitempty"`
}

// GenMessage compresses the given file and generates a
// Message instance
func GenMessage(i int, fi os.FileInfo, dir string) (*Message, error) {
	var buf bytes.Buffer
	err := util.CompressFile(filepath.Join(dir, fi.Name()), &buf)
	if err != nil {
		return nil, err
	}
	m := &Message{
		Id:           i,
		Name:         fi.Name(),
		InputPath:    filepath.Join(dir, fi.Name()),
		InputContent: buf.Bytes(),
	}
	return m, nil
}

// GenFile uncompresses the message buffer and
// generates a file
func (m *Message) GenFile(destDir string) error {
	if m.Name == "" {
		return fmt.Errorf("the file name of message %d is not given", m.Id)
	}

	// create the output dir if not exist
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		logrus.Infof("[MESSAGE] %s not exist, create it", destDir)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return err
		}
	}

	fPath := filepath.Join(destDir, m.Name)
	// uncompress the file contents
	if m.InputContent != nil && len(m.InputContent) != 0 {
		if err := util.DecompressFile(
			bytes.NewBuffer(m.InputContent), fPath); err != nil {
			return err
		}
	} else {
		logrus.Infof("[MESSAGE] the inputcontent of message %d is empty", m.Id)
		// create the empty file
		_, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
