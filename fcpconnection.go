package fcp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

type FCPConnection struct {
	host         *net.TCPAddr
	identifier   string
	isSSL        bool
	socket       *net.TCPConn
	reader       *bufio.Reader
	writer       *bufio.Writer
	log          FCPLogger
	helloMessage *fcpMessage
}

func NewFCPConnectionSimple(fcphost string) (*FCPConnection, error) {
	return NewFCPConnection(fcphost, "", false, nil)
}

func NewFCPConnection(fcphost string, identifier string, isSSL bool, logger FCPLogger) (*FCPConnection, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", fcphost)
	if err != nil {
		return nil, err

	}
	if identifier == "" {
		identifier = strconv.FormatInt(rand.Int63(), 16)
	}
	sock, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	conn := FCPConnection{
		host:       tcpAddr,
		identifier: identifier,
		isSSL:      isSSL,
		socket:     sock,
		reader:     bufio.NewReaderSize(sock, fcpMaxLineLength),
		writer:     bufio.NewWriterSize(sock, fcpMaxLineLength),
		log:        logger,
	}

	// do the client hello

	helloMessage := NewFCPMessage("ClientHello")
	helloMessage.setString("Name", identifier)
	helloMessage.setString("ExpectedVersion", "2.0")

	//	fmt.Println("Teste", conn.log)
	//	fmt.Println("Teste", logger)
	conn.SendMessage(helloMessage)
	msg, err := conn.ReadMessage()

	if err != nil {
		return &conn, err
	}

	if msg.IsMessageName("NodeHello") {
		conn.helloMessage = msg
	} else {
		return &conn, errors.New("Unecpected reply from node")
	}

	return &conn, err

	/*	ip := net.ParseIP(host)
		a := net.TCPAddr{ip, port}
		return nil*/
}

func (c *FCPConnection) SendMessage(msg *fcpMessage) (err error) {
	// TODO error handling
	_, err = c.writer.WriteString(msg.messageName)
	if err != nil {
		return err
	}
	_, err = c.writer.WriteString("\n")
	if err != nil {
		return err
	}
	if c.log != nil {
		//	fmt.Println("huhu", msg.messageName)
		//	fmt.Println("Test", c.log)
		c.log.LogSend(msg.messageName)
	}
	for key, value := range msg.params {
		c.writer.WriteString(key)
		c.writer.WriteString("=")
		c.writer.WriteString(value)
		c.writer.WriteString("\n")
		if c.log != nil {
			c.log.LogSend(fmt.Sprintf("send: %s=%s", key, value))
		}
	}
	c.writer.WriteString(msg.messageEnd)
	c.writer.WriteString("\n")
	if c.log != nil {
		c.log.LogSend(msg.messageEnd)
	}
	c.writer.Flush()
	return nil
}

func (c *FCPConnection) readLine() (string, error) {
	l, isPrefix, err := c.reader.ReadLine()
	if isPrefix {
		return string(l), errors.New("Line to long")
	}
	return string(l), err
}

func (c *FCPConnection) ReadMessage() (*fcpMessage, error) {
	msgName, err := c.readLine()
	if c.log != nil {
		c.log.LogRecv(msgName)
	}
	msg := NewFCPMessage(msgName)
	if err != nil {
		return msg, err
	}
	for {
		l, err := c.readLine()

		if c.log != nil {
			c.log.LogRecv(l)
		}

		if err != nil {
			return msg, err
		}
		s := strings.SplitN(l, "=", 2)
		switch len(s) {
		case 2:
			msg.setString(s[0], s[1])
		case 1:
			msg.messageEnd = l
			return msg, nil
		}
	}
}

func (c *FCPConnection) sendDisconnect() error {
	dmsg := NewFCPMessage("Disconnect")
	return c.SendMessage(dmsg)
}

func (c *FCPConnection) Close() {
	c.sendDisconnect()
	c.socket.Close()
}

func (c *FCPConnection) sendData(b []byte) (int, error) {
	n, err := c.writer.Write(b)
	c.writer.Flush()
	return n, err
}

func (c *FCPConnection) copyTo(src io.Reader, length int64) error {
	_, err := io.CopyN(c.writer, src, length)
	return err
}

func (c *FCPConnection) readData(p []byte) (int, error) {
	return c.reader.Read(p)
}
