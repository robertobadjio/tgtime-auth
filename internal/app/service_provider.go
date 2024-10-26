package app

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/robertobadjio/platform-common/pkg/closer"
	"github.com/robertobadjio/platform-common/pkg/db"
	"github.com/robertobadjio/platform-common/pkg/db/pg"

	"github.com/robertobadjio/tgtime-auth/internal/config"
	"github.com/robertobadjio/tgtime-auth/internal/repository/user"
	"github.com/robertobadjio/tgtime-auth/internal/repository/user/pg_db"
	"github.com/robertobadjio/tgtime-auth/internal/service/access"
	endpoints2 "github.com/robertobadjio/tgtime-auth/internal/service/access/endpoints"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/endpoints"
	"github.com/robertobadjio/tgtime-auth/internal/service/service"
	endpointsService "github.com/robertobadjio/tgtime-auth/internal/service/service/endpoints"
	transportServiceHttp "github.com/robertobadjio/tgtime-auth/internal/service/service/transport"
)

type serviceProvider struct {
	logger log.Logger

	pgConfig config.PGConfig
	db       db.Client

	grpcConfig config.GRPCConfig
	httpConfig config.HTTPConfig

	httpServiceHandler http.Handler
	endpointServiceSet endpointsService.Set
	apiServiceService  service.Service

	endpointAuthSet   endpoints.Set
	endpointAccessSet endpoints2.Set
	authService       auth.Service
	accessService     access.Service
	userRepository    user.Repository

	token config.Token
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) HTTPConfig() config.HTTPConfig {
	if sp.httpConfig == nil {
		httpConfig, err := config.NewHTTPConfig()
		if err != nil {
			_ = sp.Logger().Log("config", "http", "error", err.Error())
			os.Exit(1)
		}

		sp.httpConfig = httpConfig
	}

	return sp.httpConfig
}

func (sp *serviceProvider) HTTPServiceHandler(ctx context.Context) http.Handler {
	if sp.httpServiceHandler == nil {
		sp.httpServiceHandler = transportServiceHttp.NewHTTPHandler(sp.EndpointServiceSet(ctx))
	}

	return sp.httpServiceHandler
}

func (sp *serviceProvider) EndpointServiceSet(ctx context.Context) endpointsService.Set {
	sp.endpointServiceSet = endpointsService.NewEndpointSet(sp.APIServiceService(ctx))

	return sp.endpointServiceSet
}

func (sp *serviceProvider) APIServiceService(_ context.Context) service.Service {
	if sp.apiServiceService == nil {
		sp.apiServiceService = service.NewService()
	}

	return sp.apiServiceService
}

func (sp *serviceProvider) Token() config.Token {
	if sp.token == nil {
		token, err := config.NewToken()
		if err != nil {
			_ = sp.Logger().Log("type", "di", "service", "token", "err", err.Error())
			os.Exit(1)
		}

		sp.token = token
	}

	return sp.token
}

func (sp *serviceProvider) DB(ctx context.Context) db.Client {
	if sp.db == nil {
		cl, err := pg.New(ctx, sp.PGConfig().DSN())
		if err != nil {
			_ = sp.Logger().Log("type", "di", "service", "db client master", "err", err.Error())
			os.Exit(1)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			_ = sp.Logger().Log("type", "di", "service", "ping db client master", "err", err.Error())
			os.Exit(1)
		}
		closer.Add(cl.Close)

		sp.db = cl
	}

	return sp.db
}

func (sp *serviceProvider) PGConfig() config.PGConfig {
	if sp.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			_ = sp.Logger().Log("type", "di", "service", "pgConfig", "err", err.Error())
			os.Exit(1)
		}

		sp.pgConfig = cfg
	}
	return sp.pgConfig
}

func (sp *serviceProvider) GRPCConfig() config.GRPCConfig {
	if sp.grpcConfig == nil {
		grpcConfig, err := config.NewGRPCConfig()
		if err != nil {
			_ = sp.Logger().Log("config", "http", "error", err.Error())
			os.Exit(1)
		}

		sp.grpcConfig = grpcConfig
	}

	return sp.grpcConfig
}

func (sp *serviceProvider) UserRepository(ctx context.Context) user.Repository {
	if sp.userRepository == nil {
		sp.userRepository = pg_db.NewPgRepository(sp.DB(ctx))
	}

	return sp.userRepository
}

func (sp *serviceProvider) EndpointAuthSet(ctx context.Context) endpoints.Set {
	sp.endpointAuthSet = endpoints.NewEndpointSet(sp.AuthService(ctx))

	return sp.endpointAuthSet
}

func (sp *serviceProvider) EndpointAccessSet(ctx context.Context) endpoints2.Set {
	sp.endpointAccessSet = endpoints2.NewEndpointSet(sp.AccessService(ctx))

	return sp.endpointAccessSet
}

func (sp *serviceProvider) AuthService(ctx context.Context) auth.Service {
	if sp.authService == nil {
		sp.authService = auth.NewService(
			sp.UserRepository(ctx),
			sp.Token(),
		)
	}

	return sp.authService
}

func (sp *serviceProvider) AccessService(ctx context.Context) access.Service {
	if sp.authService == nil {
		sp.authService = auth.NewService(
			sp.UserRepository(ctx),
			sp.Token(),
		)
	}

	return sp.accessService
}

func (sp *serviceProvider) Logger() log.Logger {
	if sp.logger == nil {
		logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		sp.logger = logger
	}
	return sp.logger
}
