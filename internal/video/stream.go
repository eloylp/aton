package video

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/eloylp/aton/internal/logging"
)

type Capture struct {
	Data      []byte
	Timestamp time.Time
}

type MJPEGStreamCapturer struct {
	URL    *url.URL
	output chan *Capture
	close  chan struct{}
	logger logging.Logger
}

func NewMJPEGStreamCapturer(rawURL string, maxFrameBuffer int, logger logging.Logger) (*MJPEGStreamCapturer, error) {
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
		logger: logger,
	}, nil
}

func (m *MJPEGStreamCapturer) Start() {
	resp, err := http.Get(m.URL.String())
	if err != nil {
		m.logger.Error(fmt.Errorf("capturer: %w", err))
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
				m.logger.Error(err)
			}
			mr := multipart.NewReader(resp.Body, param["boundary"])
			for {
				p, err := mr.NextPart()
				if errors.Is(err, io.EOF) {
					m.logger.Info(fmt.Errorf("capturer: %w", err))
					break
				}
				if err != nil {
					m.logger.Error(fmt.Errorf("capturer: %w", err))
					break
				}
				data, err := ioutil.ReadAll(p)
				if err != nil {
					m.logger.Error(fmt.Errorf("capturer: %w", err))
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
