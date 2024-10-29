package app

import (
	"context"
	"net/http"

	"github.com/robertobadjio/platform-common/pkg/closer"
	"github.com/robertobadjio/platform-common/pkg/db"
	"github.com/robertobadjio/platform-common/pkg/db/pg"
	"github.com/robertobadjio/tgtime-auth/internal/config"
	"github.com/robertobadjio/tgtime-auth/internal/logger"
	"github.com/robertobadjio/tgtime-auth/internal/repository/user"
	"github.com/robertobadjio/tgtime-auth/internal/repository/user/pg_db"
	"github.com/robertobadjio/tgtime-auth/internal/service/access"
	endpointsAccess "github.com/robertobadjio/tgtime-auth/internal/service/access/endpoints"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/endpoints"
	"github.com/robertobadjio/tgtime-auth/internal/service/service"
	endpointsService "github.com/robertobadjio/tgtime-auth/internal/service/service/endpoints"
	transportServiceHttp "github.com/robertobadjio/tgtime-auth/internal/service/service/transport"
)

type serviceProvider struct {
	pgConfig config.PGConfig
	db       db.Client

	grpcConfig config.GRPCConfig
	httpConfig config.HTTPConfig

	httpServiceHandler http.Handler
	endpointServiceSet endpointsService.Set
	apiServiceService  service.Service

	endpointAuthSet endpoints.Set
	authService     auth.Service

	endpointAccessSet endpointsAccess.Set
	accessService     access.Service

	userRepository user.Repository

	token        config.Token
	promConfig   config.PromConfig
	jaegerConfig config.JaegerConfig
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) JaegerConfig() config.JaegerConfig {
	if sp.jaegerConfig == nil {
		jaegerConfig, err := config.NewJaegerConfig()
		if err != nil {
			logger.Fatal("config", "jaeger", "error", err.Error())
		}

		sp.jaegerConfig = jaegerConfig
	}

	return sp.jaegerConfig
}

func (sp *serviceProvider) PromConfig() config.PromConfig {
	if sp.promConfig == nil {
		promConfig, err := config.NewPromConfig()
		if err != nil {
			logger.Fatal("config", "http", "error", err.Error())
		}

		sp.promConfig = promConfig
	}

	return sp.promConfig
}

func (sp *serviceProvider) HTTPConfig() config.HTTPConfig {
	if sp.httpConfig == nil {
		httpConfig, err := config.NewHTTPConfig()
		if err != nil {
			logger.Fatal("config", "http", "error", err.Error())
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
			logger.Fatal("type", "di", "service", "token", "err", err.Error())
		}

		sp.token = token
	}

	return sp.token
}

func (sp *serviceProvider) DB(ctx context.Context) db.Client {
	if sp.db == nil {
		cl, err := pg.New(ctx, sp.PGConfig().DSN())
		if err != nil {
			logger.Fatal("type", "di", "service", "db client master", "err", err.Error())
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			logger.Fatal("type", "di", "service", "ping db client master", "err", err.Error())
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
			logger.Fatal("type", "di", "service", "pgConfig", "err", err.Error())
		}

		sp.pgConfig = cfg
	}
	return sp.pgConfig
}

func (sp *serviceProvider) GRPCConfig() config.GRPCConfig {
	if sp.grpcConfig == nil {
		grpcConfig, err := config.NewGRPCConfig()
		if err != nil {
			logger.Fatal("config", "http", "error", err.Error())
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

func (sp *serviceProvider) EndpointAccessSet(ctx context.Context) endpointsAccess.Set {
	sp.endpointAccessSet = endpointsAccess.NewEndpointSet(sp.AccessService(ctx))

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

func (sp *serviceProvider) AccessService(_ context.Context) access.Service {
	if sp.accessService == nil {
		sp.accessService = access.NewService(
			sp.Token(),
		)
	}

	return sp.accessService
}
