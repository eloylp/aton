package video

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type Capture struct {
	Data      []byte
	Timestamp time.Time
}

type MJPEGStreamCapturer struct {
	URL    *url.URL
	output chan *Capture
	close  chan struct{}
}

func NewMJPEGStreamCapturer(rawURL string, maxFrameBuffer int) (*MJPEGStreamCapturer, error) {
	captURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("capturer (%s): %w", rawURL, err)
	}
	if !regexp.MustCompile("http|https").MatchString(captURL.Scheme) {
		return nil, fmt.Errorf("capturer (%s): only http or https scheme supported", rawURL)
	}
	return &MJPEGStreamCapturer{
		URL:    captURL,
		output: make(chan *Capture, maxFrameBuffer),
		close:  make(chan struct{}, 1),
	}, nil
}

func (m *MJPEGStreamCapturer) Start() {
	resp, err := http.Get(m.URL.String())
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
