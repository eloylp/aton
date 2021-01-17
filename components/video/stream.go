package video

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
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
	UUIDIdent  string
	URL        *url.URL
	output     chan *Capture
	close      chan struct{}
	logger     *logrus.Logger
	status     string
	maxBackOff float64
}

func NewMJPEGCapturer(uuid, rawURL string, maxFrameBuffer int, logger *logrus.Logger) (*MJPEGCapturer, error) {
	captURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("capturer (%s): %w", rawURL, err)
	}
	if !regexp.MustCompile("https?").MatchString(captURL.Scheme) {
		return nil, fmt.Errorf("capturer (%s): only http or https scheme supported", rawURL)
	}
	var maxBackOff float64 = 16
	return &MJPEGCapturer{
		UUIDIdent:  uuid,
		URL:        captURL,
		output:     make(chan *Capture, maxFrameBuffer),
		close:      make(chan struct{}, 1),
		logger:     logger,
		status:     StatusNotRunning,
		maxBackOff: maxBackOff, // Todo think about extracting to a client.
	}, nil
}

func (m *MJPEGCapturer) UUID() string {
	return m.UUIDIdent
}

func (m *MJPEGCapturer) Start() {
	m.status = StatusRunning
	resp, err := m.connect()
	if err != nil {
		m.logger.Error(err)
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

func (m *MJPEGCapturer) connect() (resp *http.Response, err error) {
	var backoff float64 = 1
	var sleepTime time.Duration
	for {
		time.Sleep(sleepTime * time.Second)
		var req *http.Request
		req, err = http.NewRequestWithContext(context.TODO(), http.MethodGet, m.URL.String(), nil)
		if err != nil {
			return nil, err
		}
		client := &http.Client{} // TODO think about transport, timeouts and further configuration
		resp, err = client.Do(req)
		if err != nil {
			select {
			case <-m.close:
				close(m.output)
				return
			default:
				if backoff < m.maxBackOff {
					sleepTime = time.Duration(math.Exp2(backoff))
					backoff++
				}
				m.logger.Error(fmt.Errorf("capturer: %w", err))
				continue
			}
		}
		break
	}
	return
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
	select {
	case m.output <- &Capture{data, time.Now()}:
	case <-m.close:
	}
	return nil
}

func (m *MJPEGCapturer) NextOutput() (*Capture, error) {
	capt, ok := <-m.output
	if !ok {
		return nil, io.EOF
	}
	return capt, nil
}

func (m *MJPEGCapturer) Close() {
	close(m.close)
}

func (m *MJPEGCapturer) Status() string {
	return m.status
}

func (m *MJPEGCapturer) TargetURL() string {
	return m.URL.String()
}
