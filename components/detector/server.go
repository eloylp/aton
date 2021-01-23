package detector

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/eloylp/aton/components/detector/metrics"
	"github.com/eloylp/aton/components/proto"
)

type Server struct {
	service       service
	logger        *logrus.Logger
	s             *grpc.Server
	metricsServer *http.Server
	listenAddr    string
	metricsAddr   string
}

type service interface {
	proto.DetectorServer
	Shutdown()
}

func NewServer(listenAddr string, service service, metricsAddr string, m *metrics.Service, logger *logrus.Logger) *Server {
	logrusEntry := logrus.NewEntry(logger)
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	m.MustRegister(grpcMetrics)
	grpc_prometheus.EnableHandlingTimeHistogram()

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpcMetrics.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpcMetrics.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	grpc_prometheus.Register(grpcServer)
	grpcMetrics.InitializeMetrics(grpcServer)

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", m.HTTPHandler())
	s := &Server{
		service:       service,
		logger:        logger,
		listenAddr:    listenAddr,
		metricsAddr:   metricsAddr,
		metricsServer: &http.Server{Addr: metricsAddr, Handler: metricsMux},
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
		if err := gs.metricsServer.ListenAndServe(); err != http.ErrServerClosed {
			gs.logger.Errorf("metrics-server: %v", err)
		}
	}()
	return gs.s.Serve(lis)
}

func (gs *Server) watchForOSSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	recvSig := <-ch
	gs.logger.Infof("received shutdown signal %s, gracefully shutdown", recvSig.String())
	gs.Shutdown()
}

func (gs *Server) Shutdown() {
	gs.logger.Info("gracefully shutdown started.")
	const duration = 5 * time.Second
	ctx, cancl := context.WithTimeout(context.Background(), duration)
	defer cancl()
	if err := gs.metricsServer.Shutdown(ctx); err != nil {
		gs.logger.Errorf("shutdown: metrics-server: %v", err)
	}
	gs.service.Shutdown()
	gs.s.GracefulStop()
	gs.logger.Infof("stopped detector service at %s", gs.listenAddr)
	gs.logger.Infof("stopped detector metrics at %s", gs.metricsAddr)
}
