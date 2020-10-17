package video

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"time"
)

type Capture struct {
	Data      []byte
	Timestamp time.Time
}

type MJPEGStreamCapturer struct {
	URL    string
	output chan *Capture
	close  chan struct{}
}

func NewMJPEGStreamCapturer(url string, maxFrameBuffer int) *MJPEGStreamCapturer {
	return &MJPEGStreamCapturer{
		URL:    url,
		output: make(chan *Capture, maxFrameBuffer),
		close:  make(chan struct{}, 1),
	}
}

func (m *MJPEGStreamCapturer) Start() {
	resp, err := http.Get(m.URL)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
loop:
	for {
		select {
		case <-m.close:
			close(m.output)
			break loop
		default:
			_, param, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
			if err != nil {
				log.Fatal(err)
			}
			mr := multipart.NewReader(resp.Body, param["boundary"])
			for {
				p, err := mr.NextPart()
				if errors.Is(err, io.EOF) {
					log.Println(err)
					break
				}
				if err != nil {
					log.Println(err) // todo optimize this
					break
				}
				data, err := ioutil.ReadAll(p)
				if err != nil {
					log.Fatal(err)
				}
				p.Close()
				m.output <- &Capture{data, time.Now()}
			}
		}
	}
}

func (m *MJPEGStreamCapturer) Output() <-chan *Capture {
	return m.output
}

func (m *MJPEGStreamCapturer) Close() {
	close(m.close)
}
