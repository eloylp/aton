package helper

import (
	"net"
	"testing"
	"time"

	"golang.org/x/net/nettest"
)

func NetCloserServer(t *testing.T, expected int) (listener net.Listener, result chan time.Time) {
	var err error
	listener, err = nettest.NewLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	result = make(chan time.Time, expected)
	go func() {
		var accepted int
		for {
			if accepted == expected {
				close(result)
				listener.Close()
				break
			}
			c, err := listener.Accept()
			if err != nil {
				close(result)
				t.Log(err)
				break
			}
			result <- time.Now()
			accepted++
			if c.Close() != nil {
				t.Log(err)
			}
		}
	}()
	return
}
