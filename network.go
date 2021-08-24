package logkit

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type NetworkTarget struct {
	Network string
	Address string

	entries chan *Entry
	conn    net.Conn
	close   chan bool
}

// NewNetworkTarget return the NetworkTarget instance
func NewNetworkTarget() *NetworkTarget {
	return &NetworkTarget{
		close: make(chan bool),
	}
}

// Open init NetworkTarget property, connect to remote and keep alive
// Start a goroutine listen message, if message exist, send message
// To remote.
func (t *NetworkTarget) Open(w io.Writer) error {
	if t.Network == "" {
		return errors.New("NetworkTarget.Network must be set. ")
	}

	if t.Address == "" {
		return errors.New("NetworkTarget.Address must be set. ")
	}

	t.entries = make(chan *Entry, 1)
	t.conn = nil

	err := t.connect()
	if err != nil {
		return err
	}

	// start goroutine to send message to remote
	go t.send(w)

	return nil
}

func (t *NetworkTarget) Process(entry *Entry) {
	if entry == nil {
		close(t.entries)
	} else {
		t.entries <- entry
	}
}

func (t *NetworkTarget) Close() {
	t.close <- true
}

func (t *NetworkTarget) connect() error {
	if t.conn != nil {
		_ = t.conn.Close()
		t.conn = nil
	}
	conn, err := net.Dial(t.Network, t.Address)
	if err != nil {
		return err
	}
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetKeepAlive(true)
	}
	t.conn = conn

	return nil
}

func (t *NetworkTarget) send(errWriter io.Writer) {
	for {
		entry, ok := <-t.entries
		if !ok {
			if t.conn != nil {
				_ = t.conn.Close()
				t.conn = nil
			}
			t.close <- true
			break
		}

		if _, err := t.write(entry.String() + "\n"); err != nil {
			_, _ = fmt.Fprintf(errWriter, "NetworkTarget write error: %v\n", err)
		}
	}
}

func (t *NetworkTarget) write(message string) (int, error) {
	n, err := t.conn.Write([]byte(message))
	return n, err
}
