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

const (
	StatusNotRunning = "NOT_RUNNING"
	StatusRunning    = "RUNNING"
)

type Capture struct {
	Data      []byte
	Timestamp time.Time
}

type MJPEGCapturer struct {
	URL    *url.URL
	output chan *Capture
	close  chan struct{}
	logger logging.Logger
	status string
}

func NewMJPEGCapturer(rawURL string, maxFrameBuffer int, logger logging.Logger) (*MJPEGCapturer, error) {
	captURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("capturer (%s): %w", rawURL, err)
	}
	if !regexp.MustCompile("http|https").MatchString(captURL.Scheme) {
		return nil, fmt.Errorf("capturer (%s): only http or https scheme supported", rawURL)
	}
	return &MJPEGCapturer{
		URL:    captURL,
		output: make(chan *Capture, maxFrameBuffer),
		close:  make(chan struct{}, 1),
		logger: logger,
		status: StatusNotRunning,
	}, nil
}

func (m *MJPEGCapturer) Start() {
	m.status = StatusRunning
	resp, err := http.Get(m.URL.String())
	if err != nil {
		m.logger.Error(fmt.Errorf("capturer: %w", err))
		return
	}
	defer resp.Body.Close()
mainLoop:
	for {
		select {
		case <-m.close:
			close(m.output)
			break mainLoop
		default:
			_, param, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
			if err != nil {
				m.logger.Error(err)
				break mainLoop
			}
			mr := multipart.NewReader(resp.Body, param["boundary"])
			for {
				if err := m.processNextPart(mr); err != nil {
					m.logger.Error(err)
					break
				}
				select {
				case <-m.close:
					close(m.output)
					break mainLoop
				default:
					continue
				}
			}
		}
	}
}

func (m *MJPEGCapturer) processNextPart(mr *multipart.Reader) error {
	p, err := mr.NextPart()
	if errors.Is(err, io.EOF) {
		return fmt.Errorf("capturer: %w", err)
	}
	if err != nil {
		return fmt.Errorf("capturer: %w", err)
	}
	data, err := ioutil.ReadAll(p)
	if err != nil {
		return fmt.Errorf("capturer: %w", err)
	}
	m.output <- &Capture{data, time.Now()}
	return nil
}

func (m *MJPEGCapturer) Output() <-chan *Capture {
	return m.output
}

func (m *MJPEGCapturer) Close() {
	close(m.close)
}

func (m *MJPEGCapturer) Status() string {
	return m.status
}
