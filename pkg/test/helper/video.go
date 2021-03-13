package helper

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"testing"
	"time"
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
	}, 25))
	return httptest.NewServer(mux)
}

func streamHandler(t *testing.T, fg frameGenerator, fps int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mp := multipart.NewWriter(w)
		defer mp.Close()
		if err := mp.SetBoundary("mjpeg"); err != nil {
			t.Fatal(err)
		}
		w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+mp.Boundary())
		timePerFrame := time.Duration(float64(1) / float64(fps) * 1e9) // nanoseconds per frame
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
			now := time.Now()
			if _, err := mw.Write(frame); err != nil {
				t.Log(err)
				return
			}
			elapsed := time.Since(now)
			duration := timePerFrame - elapsed
			time.Sleep(duration)
		}
	}
}

func ReplayedVideoStream(t *testing.T, picturesPaths []string, servingPath string, times, fps int) *httptest.Server {
	t.Helper()
	frames := make([][]byte, len(picturesPaths))
	for i := 0; i < len(picturesPaths); i++ {
		frames[i] = ReadFile(t, picturesPaths[i])
	}
	rplStream := make(chan []byte, len(frames)*times)
	for i := 0; i < times; i++ {
		for f := 0; f < len(frames); f++ {
			rplStream <- frames[f]
		}
	}
	close(rplStream)
	mux := http.NewServeMux()
	mux.HandleFunc(servingPath, streamHandler(t, func() ([]byte, error) {
		frame, ok := <-rplStream
		if !ok {
			return nil, io.EOF
		}
		return frame, nil
	}, fps))
	return httptest.NewServer(mux)
}
