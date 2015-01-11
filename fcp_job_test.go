package fcp

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var stop bool
var URI string
var result []byte

func cb(msg *fcpMessage, length int64, data io.Reader, err error) {
	fmt.Println("Huhu callback", msg.messageName)
	switch msg.messageName {
	case "PutSuccessful":
		URI = msg.params["URI"]
		stop = true
	case "PutFailed":
		stop = true
	case "ProtocolError":
		stop = true
	}
}

func cb_data(msg *fcpMessage, length int64, data io.Reader, err error) {
	fmt.Println("Huhu callback data", msg.messageName)
	switch msg.messageName {
	case "AllData":
		result = make([]byte, length)
		io.ReadFull(data, result)
		stop = true
	case "GetFailed":
		stop = true
	case "ProtocolError":
		stop = true
	}
}

func TestFCPConnectionAsync(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	stop = false
	var buf bytes.Buffer
	buf.WriteString("testdata trätrööt ")

	buf.WriteString(strconv.FormatInt(rand.Int63(), 16))
	testData := buf.Bytes()
	//fmt.Println(buf.String())

	// connect
	cr, err := newFCPConnectionRunnerSimple("127.0.0.1:9481", cb)
	if err != nil {
		t.Error("Pfehler", err)
		return
	}

	// send insert
	insertMsg := newFCPMessage("ClientPut")
	insertMsg.setAutoIdentifier()
	insertMsg.setString("URI", "CHK@")
	insertMsg.setInteger("Verbosity", -1)
	insertMsg.setBoolean("Global", false)
	insertMsg.setString("Persistence", "connection")
	insertMsg.setString("UploadFrom", "direct")
	insertMsg.setInteger("DataLength", len(testData))
	insertMsg.setEndMarker(true)

	cr.sendMessageByteData(insertMsg, testData)

	for {
		if stop {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// TODO disconnect

	stop = false

	// connect again
	cr, err = newFCPConnectionRunnerSimple("127.0.0.1:9481", cb_data)
	if err != nil {
		t.Error("Pfehler", err)
	}

	// send fetch request

	fetchMsg := newFCPMessage("ClientGet")
	fetchMsg.setAutoIdentifier()
	fetchMsg.setString("URI", URI)
	fetchMsg.setInteger("Verbosity", -1)
	fetchMsg.setBoolean("Global", false)
	fetchMsg.setString("Persistence", "connection")

	// TODO check return values
	cr.sendMessage(fetchMsg)

	for {
		if stop {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("sent   : ", string(testData))
	fmt.Println("recived: ", string(result))

}
