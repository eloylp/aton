package detector

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/eloylp/aton/internal/logging"
)

type GRPCServer struct {
	service    DetectorServer
	logger     logging.Logger
	s          *grpc.Server
	listenAddr string
}

func NewGRPCServer(service DetectorServer, logger logging.Logger, listenAddr string) *GRPCServer {
	s := &GRPCServer{
		service:    service,
		logger:     logger,
		listenAddr: listenAddr,
		s:          grpc.NewServer(),
	}
	RegisterDetectorServer(s.s, s.service)
	return s
}

func (gs *GRPCServer) Start() error {
	gs.logger.Infof("starting detector service at %s", gs.listenAddr)
	lis, err := net.Listen("tcp", gs.listenAddr)
	if err != nil {
		return err
	}
	go gs.watchForOSSignals()
	return gs.s.Serve(lis)
}

func (gs *GRPCServer) watchForOSSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	recvSig := <-ch
	gs.logger.Infof("received shutdown signal %q, gracefully shutdown", recvSig.String())
	gs.Shutdown()
}

func (gs *GRPCServer) Shutdown() {
	gs.s.GracefulStop()
	gs.logger.Infof("stopped detector service at %s", gs.listenAddr)
}
