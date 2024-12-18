package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/oklog/oklog/pkg/group"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robertobadjio/platform-common/pkg/closer"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/robertobadjio/tgtime-auth/internal/interceptor"
	"github.com/robertobadjio/tgtime-auth/internal/logger"
	"github.com/robertobadjio/tgtime-auth/internal/metric"
	transportAccess "github.com/robertobadjio/tgtime-auth/internal/service/access/transport"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/transport"
	transportServiceHttp "github.com/robertobadjio/tgtime-auth/internal/service/service/transport"
	"github.com/robertobadjio/tgtime-auth/pkg/api/access_v1"
	"github.com/robertobadjio/tgtime-auth/pkg/api/auth_v1"
)

const serviceName = "tgtime-auth-service"
const promMetricsPath = "/metrics"

// App ???
type App struct {
	serviceProvider *serviceProvider
	apiGateway      group.Group
}

// NewApp ???
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run ???
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runAPIGateway()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initAPIGateway,
		a.initPrometheus,
		a.initTracing,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			logger.Fatal(
				"init", "deps",
				"error", err.Error(),
			)
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initAPIGateway(ctx context.Context) error {
	var g group.Group
	{
		httpListener, err := net.Listen("tcp", a.serviceProvider.HTTPConfig().Address())
		if err != nil {
			return err
		}
		g.Add(func() error {
			logger.Info(
				"transport", "HTTP",
				"component", "API",
				"addr", a.serviceProvider.HTTPConfig().Address(),
			)

			sm := http.NewServeMux()
			sm.Handle(
				transportServiceHttp.ServiceStatus,
				a.serviceProvider.HTTPServiceHandler(ctx),
			)

			srv := &http.Server{
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				Handler:      sm,
			}
			return srv.Serve(httpListener)
		}, func(err error) {
			logger.Error("transport", "HTTP", "component", "API", "during", "Listen", "err", err.Error())
			_ = httpListener.Close()
		})
	}
	{
		grpcListener, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
		if err != nil {
			return err
		}
		g.Add(func() error {
			logger.Info("transport", "GRPC", "component", "API", "addr", a.serviceProvider.GRPCConfig().Address())

			baseServer := grpc.NewServer(
				grpc.UnaryInterceptor(
					grpcMiddleware.ChainUnaryServer(
						kitgrpc.Interceptor,
						interceptor.ErrorCodesInterceptor,
						interceptor.LogInterceptor,
						interceptor.MetricsInterceptor,
						interceptor.ServerTracingInterceptor,
					),
				),
			)

			reflection.Register(baseServer)

			auth_v1.RegisterAuthV1Server(
				baseServer,
				transport.NewGRPCServer(a.serviceProvider.EndpointAuthSet(ctx)),
			)
			access_v1.RegisterAccessV1Server(
				baseServer,
				transportAccess.NewGRPCServer(a.serviceProvider.EndpointAccessSet(ctx)),
			)

			return baseServer.Serve(grpcListener)
		}, func(err error) {
			logger.Error("transport", "GRPC", "component", "API", "during", "Listen", "err", err.Error())
			_ = grpcListener.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	a.apiGateway = g

	return nil
}

func (a *App) initPrometheus(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		return err
	}

	httpListener, err := net.Listen("tcp", a.serviceProvider.PromConfig().Address())
	if err != nil {
		return err
	}

	a.apiGateway.Add(func() error {
		logger.Info(
			"transport", "HTTP",
			"component", "prometheus",
			"addr", a.serviceProvider.PromConfig().Address(),
		)

		sm := http.NewServeMux()
		sm.Handle(promMetricsPath, promhttp.Handler())

		srv := &http.Server{
			Handler:      sm,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		return srv.Serve(httpListener)
	}, func(err error) {
		logger.Error("transport", "HTTP", "component", "prometheus", "during", "listen", "err", err.Error())
		_ = httpListener.Close()
	})

	return nil
}

func (a *App) initTracing(_ context.Context) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: a.serviceProvider.JaegerConfig().Address(),
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runAPIGateway() error {
	return a.apiGateway.Run()
}
