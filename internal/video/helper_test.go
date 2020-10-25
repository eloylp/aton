package video_test

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"testing"
	"time"

	"golang.org/x/net/nettest"
)

func videoStream(t *testing.T, picturesPaths []string, servingPath string) *httptest.Server {
	t.Helper()
	pictures := make([][]byte, len(picturesPaths))
	for i := 0; i < len(picturesPaths); i++ {
		data, err := ioutil.ReadFile(picturesPaths[i])
		if err != nil {
			t.Fatal(err)
		}
		pictures[i] = data
	}
	mux := http.NewServeMux()
	mux.HandleFunc(servingPath, func(w http.ResponseWriter, r *http.Request) {
		mp := multipart.NewWriter(w)
		defer mp.Close()
		if err := mp.SetBoundary("mjpeg"); err != nil {
			t.Fatal(err)
		}
		w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+mp.Boundary())
		for i := 0; i < len(pictures); i++ {
			h := textproto.MIMEHeader{}
			pictureSize := len(pictures[i])
			h.Add("Content-Type", "image/jpeg")
			h.Add("Content-Length", strconv.Itoa(pictureSize))
			mw, err := mp.CreatePart(h)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := mw.Write(pictures[i]); err != nil {
				t.Log(err)
				return
			}
		}
	})
	return httptest.NewServer(mux)
}

func readFile(t *testing.T, file string) []byte {
	t.Helper()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func netConnection(t *testing.T, expected int) (net.Listener, chan time.Time) {
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
				break
			}
			c, err := s.Accept()
			if err != nil {
				close(m)
				t.Log(err)
				break
			}
			accepted++
			m <- time.Now()
			if c.Close() != nil {
				t.Log(err)
			}
		}
	}()
	return s, m
}
