package helper

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"testing"
)

type frameGenerator func() ([]byte, error)

func VideoStream(t *testing.T, picturesPaths []string, servingPath string) *httptest.Server {
	t.Helper()
	frames := make(chan []byte, len(picturesPaths))
	for i := 0; i < len(picturesPaths); i++ {
		frames <- ReadFile(t, picturesPaths[i])
	}
	close(frames)
	mux := http.NewServeMux()
	mux.HandleFunc(servingPath, streamHandler(t, func() ([]byte, error) {
		frame, ok := <-frames
		if !ok {
			return nil, io.EOF
		}
		return frame, nil
	}))
	return httptest.NewServer(mux)
}

func streamHandler(t *testing.T, fg frameGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mp := multipart.NewWriter(w)
		defer mp.Close()
		if err := mp.SetBoundary("mjpeg"); err != nil {
			t.Fatal(err)
		}
		w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+mp.Boundary())
		for {
			frame, err := fg()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			h := textproto.MIMEHeader{}
			h.Add("Content-Type", "image/jpeg")
			h.Add("Content-Length", strconv.Itoa(len(frame)))
			mw, err := mp.CreatePart(h)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := mw.Write(frame); err != nil {
				t.Log(err)
				return
			}
		}
	}
}
