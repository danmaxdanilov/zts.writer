package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/danmaxdanilov/zts.shared/pkg/interceptors"
	"github.com/danmaxdanilov/zts.shared/pkg/logger"
	"github.com/danmaxdanilov/zts.shared/pkg/tracing"
	"github.com/danmaxdanilov/zts.writer/config"
	"github.com/danmaxdanilov/zts.writer/internal/metrics"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type server struct {
	log       logger.Logger
	cfg       *config.Config
	v         *validator.Validate
	kafkaConn *kafka.Conn
	im        interceptors.InterceptorManager
	pgConn    *pgxpool.Pool
	metrics   *metrics.WriterServiceMetrics
	//ps        *service.ProductService
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	defer cancel()

	s.runHealthCheck(ctx)
	s.runMetrics(cancel)

	if s.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
	}

	<-ctx.Done()

	return nil
}
