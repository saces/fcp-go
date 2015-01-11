package fcp

import (
	"fmt"
)

type FCPLogger interface {
	LogRecv(s string)
	LogRecvData(i int)
	LogSend(s string)
	LogSendData(i int)
}

type SimpleLogger struct{}

func (l *SimpleLogger) LogRecv(s string) {
	fmt.Println("recv:", s)
}
func (l *SimpleLogger) LogRecvData(i int) {
	fmt.Printf("recv: %d Bytes\n", i)
}
func (l *SimpleLogger) LogSend(s string) {
	fmt.Println("send:", s)
}
func (l *SimpleLogger) LogSendData(i int) {
	fmt.Printf("send: %d Bytes\n", i)
}

var _ FCPLogger = (*SimpleLogger)(nil)
