package builder

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/egress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/services"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/egress/cache"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/egress/database"
	databaseRepository "github.com/bhupendra-dudhwal/sso-gateway/internal/egress/repository/database"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/ingress/handler"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	"github.com/bhupendra-dudhwal/sso-gateway/pkg/logger"
	"github.com/fasthttp/router"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	httpClient "github.com/bhupendra-dudhwal/sso-gateway/internal/egress/http"
)

type appBuilder struct {
	ctx context.Context

	config *models.Config
	// logger ports.Logger

	// // Client
	// dbClient    *gorm.DB
	// cacheClient *redis.Client

	// httpClient egress.HttpClientPorts

	handler *router.Router
	server  *fasthttp.Server

	// repositories
	ingressRepository ingress.Repository
	egressRepository  egress.Repository
	repository        ports.Repository

	// Services
	// healthService     ingress.HealthServicePorts
	// tokenService      ingress.TokenServicePorts
	// roleService       ingress.RoleServicePorts
	// permissionService ingress.PermissionServicePorts
	// userService       ingress.UserServicePorts
	// authService       ingress.AuthServicePorts
}

func NewAppBuilder(ctx context.Context) *appBuilder {
	return &appBuilder{
		ctx:    ctx,
		server: &fasthttp.Server{},
	}
}

func (a *appBuilder) SetConfig() *appBuilder {
	rootDir, err := utils.FindProjectRoot()
	if err != nil {
		log.Fatalf("[config] Failed to find project root: %v", err)
	}

	configFilePath := filepath.Join(rootDir, "config", "config.yaml")

	if !utils.FileExists(configFilePath) {
		log.Fatalf("[config] Config file not found at: %s", configFilePath)
	}

	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("[config] Failed to read config file: %v", err)
	}

	cnf := models.Config{}
	if err := yaml.Unmarshal(configBytes, &cnf); err != nil {
		log.Fatalf("[config] Failed to unmarshal YAML: %v", err)
	}

	if err := cnf.Validate(); err != nil {
		log.Fatalf("[config] Invalid config: %v", err)
	}

	a.config = &cnf

	return a
}

func (a *appBuilder) SetLogger() *appBuilder {
	a.repository.Logger = logger.NewLogger(a.config.Logger, a.config.App.Server.Environment)
	return a
}

// Call from repository
func (a *appBuilder) setDatabase() *gorm.DB {
	db, err := database.NewDatabase(a.config.Database, a.repository.Logger).Connect()
	if err != nil {
		a.repository.Logger.Error("Database connection error", zap.Error(err))
		os.Exit(1)
	}

	return db
}

// Call from repository
func (a *appBuilder) setCache() *redis.Client {
	db, err := cache.NewCache(a.config.Cache, a.repository.Logger).Connect(a.ctx)
	if err != nil {
		a.repository.Logger.Error("Cache connection error", zap.Error(err))
		os.Exit(1)
	}

	return db
}

func (a *appBuilder) SetDatabaseRepositories() *appBuilder {
	client := a.setDatabase()

	a.egressRepository.Role = databaseRepository.NewRoleRepository(client)
	a.egressRepository.User = databaseRepository.NewUserRepository(client)
	a.egressRepository.LoginHistory = databaseRepository.NewloginHistoryRepository(client)
	a.egressRepository.Permission = databaseRepository.NewPermissionRepository(client)

	return a
}

func (a *appBuilder) SetServices() *appBuilder {
	a.ingressRepository.Health = services.NewHealthService(a.config, a.repository.Logger)
	a.ingressRepository.Token = services.NewTokenService(a.config, a.repository.Logger)
	a.ingressRepository.Auth = services.NewAuthService(
		a.config, a.repository.Logger, a.egressRepository.User, a.egressRepository.Role, a.ingressRepository.Token,
		a.egressRepository.LoginHistory,
	)
	a.ingressRepository.Role = services.NewRoleService(a.config, a.repository.Logger, a.egressRepository)
	a.ingressRepository.Permission = services.NewPermissionService(a.config, a.repository.Logger, a.egressRepository)
	a.ingressRepository.User = services.NewUserService(a.config, a.repository.Logger)

	return a
}

func (a *appBuilder) SetHandler() *appBuilder {
	routes, handlerObj := handler.NewHandler()

	handlerObj.SetHealthHandler(a.ingressRepository.Health)
	handlerObj.SetRoleHandler(a.ingressRepository.Role)
	handlerObj.SetPermissionHandler(a.ingressRepository.Permission)
	handlerObj.SetUserHandler(a.ingressRepository.User)

	a.handler = routes

	return a
}

func (a *appBuilder) SetHttpClient() *appBuilder {
	client, err := httpClient.NewHttpClient(a.config.HttpClient)
	if err != nil {
		a.repository.Logger.Error("http client error", zap.Error(err))
		os.Exit(1)
	}
	a.egressRepository.HttpClient = client

	return a
}

func (a *appBuilder) Build() (ports.Logger, *fasthttp.Server, int) {
	a.server.Handler = a.handler.Handler
	return a.repository.Logger, a.server, a.config.App.Server.Port
}
