package helper

import (
	"context"
	"net"
	"testing"
	"time"
)

func TryConnectTo(t *testing.T, addr string, maxWait time.Duration) {
	ctx, cancl := context.WithTimeout(context.Background(), maxWait)
	defer cancl()
	var con net.Conn
	var conErr error
mainLoop:
	for {
		select {
		case <-ctx.Done():
			t.Errorf("TryConnectTo(): %s:%s", ctx.Err(), conErr)
		default:
			con, conErr = net.Dial("tcp", addr)
			if conErr == nil {
				con.Close()
				break mainLoop
			}
		}
	}
}
