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

// simple & synchronous & blocking
func TestFCPConnection(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	var buf bytes.Buffer
	buf.WriteString("testdata trätrööt ")
	buf.WriteString(strconv.FormatInt(rand.Int63(), 16))
	testData := buf.Bytes()
	//fmt.Println(buf.String())

	// connect
	conn, err := newFCPConnectionSimple("127.0.0.1:9481")
	if err != nil {
		t.Error("Pfehler", err)
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

	// TODO check return values
	conn.sendMessage(insertMsg)
	conn.sendData(testData)

	var URI string // URI = insert result
	var msg *fcpMessage

OuterFor1:
	for {
		msg, err = conn.readMessage()
		if err != nil {
			t.Error("Pfehler", err)
		}
		switch msg.messageName {
		case "ProtocolError":
			t.Error("Pfehler")
		case "PutFailed":
			t.Error("Pfehler")
		case "PutSuccessful":
			URI = msg.params["URI"]
			break OuterFor1
		}
	}

	err = conn.sendDisconnect()
	if err != nil {
		t.Error("Pfehler", err)
	}
	_, err = conn.readMessage()
	if err != nil && err != io.EOF {
		t.Error("Pfehler", err)
	}

	//

	// connect again
	conn, err = newFCPConnectionSimple("127.0.0.1:9481")
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
	conn.sendMessage(fetchMsg)

	var size int
OuterFor2:
	for {
		msg, err = conn.readMessage()
		if err != nil {
			t.Error("Pfehler", err)
		}
		switch msg.messageName {
		case "ProtocolError":
			t.Error("Pfehler")
		case "GetFailed":
			t.Error("Pfehler")
		case "AllData":
			size, err = msg.getInteger("DataLength")
			if err != nil {
				t.Error("Pfehler", err)
			}
			break OuterFor2
		}
	}

	result := make([]byte, size)
	conn.readData(result)

	fmt.Println("sent   :", string(testData))
	fmt.Println("recived:", string(result))

	err = conn.sendDisconnect()
	if err != nil {
		t.Error("Pfehler", err)
	}
	_, err = conn.readMessage()
	if err != nil && err != io.EOF {
		t.Error("Pfehler", err)
	}
}
