package helper

import (
	"net"
	"testing"
	"time"

	"golang.org/x/net/nettest"
)

func NetCloserServer(t *testing.T, expected int) (net.Listener, chan time.Time) {
	s, err := nettest.NewLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	m := make(chan time.Time, expected)
	go func() {
		var accepted int
		for {
			if accepted == expected {
				close(m)
				s.Close()
				break
			}
			c, err := s.Accept()
			if err != nil {
				close(m)
				t.Log(err)
				break
			}
			m <- time.Now()
			accepted++
			if c.Close() != nil {
				t.Log(err)
			}
		}
	}()
	return s, m
}
