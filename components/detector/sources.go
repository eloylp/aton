package detector

import (
	"github.com/eloylp/aton/components/video"
)

type Capturer interface {
	Start()
	NextOutput() (*video.Capture, error)
	Close()
	Status() string
	UUID() string
	TargetURL() string
}
