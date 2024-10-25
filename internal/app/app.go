package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/robertobadjio/platform-common/pkg/closer"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/transport"
	"github.com/robertobadjio/tgtime-auth/pkg/api/auth_v1"
)

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
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
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
		grpcListener, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
		if err != nil {
			return err
		}
		g.Add(func() error {
			_ = a.serviceProvider.Logger().Log("transport", "GRPC", "addr", a.serviceProvider.GRPCConfig().Address())

			baseServer := grpc.NewServer(
				grpc.UnaryInterceptor(kitgrpc.Interceptor),
			)

			reflection.Register(baseServer)

			auth_v1.RegisterAuthV1Server(baseServer, transport.NewGRPCServer(a.serviceProvider.EndpointAuthSet(ctx)))
			//access_v1.RegisterAccessV1Server(baseServer, NewGRPCServer(a.serviceProvider.EndpointAccessSet(ctx)))

			return baseServer.Serve(grpcListener)
		}, func(error) {
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

func (a *App) runAPIGateway() error {
	return a.apiGateway.Run()
}
