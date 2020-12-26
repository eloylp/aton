package grpc

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/eloylp/aton/internal/logging"
	"github.com/eloylp/aton/internal/proto"
)

type Server struct {
	service       proto.DetectorServer
	logger        logging.Logger
	s             *grpc.Server
	metricsServer *http.Server
	listenAddr    string
	metricsAddr   string
}

func NewServer(listenAddr string, service proto.DetectorServer, metricsAddr string, logger logging.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	grpc_prometheus.Register(grpcServer)
	http.Handle("/metrics", promhttp.Handler())
	s := &Server{
		service:       service,
		logger:        logger,
		listenAddr:    listenAddr,
		metricsAddr:   metricsAddr,
		metricsServer: &http.Server{Addr: metricsAddr},
		s:             grpcServer,
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
	gs.logger.Infof("starting detector metrics at %s", gs.metricsAddr)
	go func() {
		if err := http.ListenAndServe(gs.metricsAddr, nil); err != http.ErrServerClosed {
			gs.logger.Errorf("metrics-server: %v", err)
		}
	}()
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
	const duration = 5 * time.Second
	ctx, cancl := context.WithTimeout(context.Background(), duration)
	defer cancl()
	if err := gs.metricsServer.Shutdown(ctx); err != nil {
		gs.logger.Errorf("shutdown: metrics-server: %v", err)
	}
	gs.s.GracefulStop()
	gs.logger.Infof("stopped detector service at %s", gs.listenAddr)
	gs.logger.Infof("stopped detector metrics at %s", gs.metricsAddr)
}
