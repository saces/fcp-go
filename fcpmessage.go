package fcp

import (
	"math/rand"
	"strconv"
	"strings"
)

type fcpMessage struct {
	messageName string
	params      map[string]string
	messageEnd  string
}

func NewFCPMessage(msgName string) *fcpMessage {
	return &fcpMessage{
		messageName: msgName,
		params:      make(map[string]string),
		messageEnd:  msgEndMarker,
	}
}

func (m *fcpMessage) setEndMarker(hasData bool) {
	if hasData {
		m.messageEnd = msgEndMarkerData
	} else {
		m.messageEnd = msgEndMarker
	}
}

func (m *fcpMessage) hasData() bool {
	return m.messageEnd == msgEndMarkerData
}

func (m *fcpMessage) setItem(key, value string) {
	m.params[key] = value
}

func (m *fcpMessage) setString(key, value string) {
	m.setItem(key, value)
}

func (m *fcpMessage) setBoolean(key string, value bool) {
	if value {
		m.setItem(key, "true")

	} else {
		m.setItem(key, "false")
	}
}

func (m *fcpMessage) setInteger(key string, value int) {
	m.setItem(key, strconv.Itoa(value))
}

func (m *fcpMessage) IsMessageName(name string) bool {
	return strings.EqualFold(m.messageName, name)
}

func (m *fcpMessage) setAutoIdentifier() {
	m.setItem("Identifier", strconv.FormatInt(rand.Int63(), 16))
}

func (m *fcpMessage) getInteger(name string) (int, error) {
	return strconv.Atoi(m.params[name])
}

func (m *fcpMessage) getInt64(name string) (int64, error) {
	return strconv.ParseInt(m.params[name], 0, 64)
}

func (m *fcpMessage) GetString(name string) string {
	return m.params[name]
}
