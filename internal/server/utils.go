package server

import (
	"context"
	"net/http"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	stackSize = 1 << 10 // 1KB
)

func (s *server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()

	health.AddLivenessCheck(s.cfg.ServiceName, healthcheck.AsyncWithContext(ctx, func() error {
		return nil
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck("Simple check", healthcheck.Async(func() error {
		return nil
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	// health.AddReadinessCheck(constants.Postgres, healthcheck.AsyncWithContext(ctx, func() error {
	// 	return s.pgConn.Ping(ctx)
	// }, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	// health.AddReadinessCheck(constants.Kafka, healthcheck.AsyncWithContext(ctx, func() error {
	// 	_, err := s.kafkaConn.Brokers()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	go func() {
		s.log.Infof("Writer microservice Kubernetes probes listening on port: %s", s.cfg.Probes.Port)
		if err := http.ListenAndServe(s.cfg.Probes.Port, health); err != nil {
			s.log.WarnMsg("ListenAndServe", err)
		}
	}()
}

func (s *server) runMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	go func() {
		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:         stackSize,
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
		metricsServer.GET(s.cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		s.log.Infof("Metrics server is running on port: %s", s.cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(s.cfg.Probes.PrometheusPort); err != nil {
			s.log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()
}
