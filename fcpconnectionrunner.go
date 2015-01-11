package fcp

import (
	"bytes"
	"fmt"
	"io"
)

type callBack func(msg *fcpMessage, length int64, data io.Reader, err error)

type sendX struct {
	msg    *fcpMessage
	length int64
	data   io.Reader
}

type recvX struct {
	msg    *fcpMessage
	length int64
	data   io.Reader
}

type fcpConnectionRunner struct {
	c         *FCPConnection
	cb        callBack
	lastError *error
	sendQ     chan sendX
	recvQ     chan recvX
}

func newFCPConnectionRunnerSimple(fcphost string, cb callBack) (*fcpConnectionRunner, error) {
	return newFCPConnectionRunner(fcphost, "", false, cb)
}

func newFCPConnectionRunner(fcphost string, identifier string, isSSL bool, cb callBack) (cr *fcpConnectionRunner, err error) {

	cr = &fcpConnectionRunner{}
	cr.cb = cb

	// make channels
	cr.sendQ = make(chan sendX)
	cr.recvQ = make(chan recvX)

	//sendStop := make(chan int)
	//recvStop := make(chan int)

	// setup connection
	cr.c, err = NewFCPConnection(fcphost, identifier, isSSL, nil)
	if err != nil {
		return cr, err
	}
	// start it
	go cr.doSender()
	go cr.doReciver()

	return cr, err
}

func (cr *fcpConnectionRunner) doSender() {
	for {
		s := <-cr.sendQ
		cr.c.SendMessage(s.msg)
		if s.data != nil {
			cr.c.copyTo(s.data, s.length)
		}
	}
}

func (cr *fcpConnectionRunner) doReciver() {
	for {
		msg, err := cr.c.ReadMessage()
		if !msg.hasData() {
			cr.cb(msg, 0, nil, err)
		} else {
			l, err := msg.getInt64("DataLength")
			r := &io.LimitedReader{cr.c.reader, l}
			cr.cb(msg, l, r, err)
			fmt.Println("Callback returned, read", r.N)
		}
	}
}

func (cr *fcpConnectionRunner) sendMessage(msg *fcpMessage) {
	cr.sendMessageData(msg, 0, nil)
}

func (cr *fcpConnectionRunner) sendMessageByteData(msg *fcpMessage, b []byte) {
	cr.sendMessageData(msg, int64(len(b)), bytes.NewReader(b))
}

func (cr *fcpConnectionRunner) sendMessageData(msg *fcpMessage, length int64, r io.Reader) {
	m := sendX{msg, length, r}
	cr.sendQ <- m
}
