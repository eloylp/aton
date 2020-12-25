package grpc

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/eloylp/aton/internal/logging"
	"github.com/eloylp/aton/internal/proto"
)

type Server struct {
	service     proto.DetectorServer
	logger      logging.Logger
	s           *grpc.Server
	listenAddr  string
	metricsAddr string
}

func NewServer(listenAddr string, service proto.DetectorServer, metricsAddr string, logger logging.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	grpc_prometheus.Register(grpcServer)
	http.Handle("/metrics", promhttp.Handler())
	s := &Server{
		service:     service,
		logger:      logger,
		listenAddr:  listenAddr,
		metricsAddr: metricsAddr,
		s:           grpcServer,
	}
	proto.RegisterDetectorServer(grpcServer, service)
	return s
}

func (gs *Server) Start() error {
	gs.logger.Infof("starting detector service at %s", gs.listenAddr)
	lis, err := net.Listen("tcp", gs.listenAddr)
	if err != nil {
		return err
	}
	go gs.watchForOSSignals()
	go http.ListenAndServe(gs.metricsAddr, nil)
	return gs.s.Serve(lis)
}

func (gs *Server) watchForOSSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	recvSig := <-ch
	gs.logger.Infof("received shutdown signal %q, gracefully shutdown", recvSig.String())
	gs.Shutdown()
}

func (gs *Server) Shutdown() {
	gs.s.GracefulStop()
	gs.logger.Infof("stopped detector service at %s", gs.listenAddr)
}
