package builder

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/services"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/egress/cache"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/egress/database"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/ingress/handler"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	"github.com/bhupendra-dudhwal/sso-gateway/pkg/logger"
	"github.com/fasthttp/router"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type appBuilder struct {
	ctx context.Context

	config *models.Config
	logger ports.Logger

	dbClient    *gorm.DB
	cacheClient *redis.Client

	handler *router.Router
	server  *fasthttp.Server

	healthService ingress.HealthServicePorts
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
	a.logger = logger.NewLogger(a.config.Logger, a.config.App.Server.Environment)
	return a
}

func (a *appBuilder) SetDatabase() *appBuilder {
	db, err := database.NewDatabase(a.config.Database, a.logger).Connect()
	if err != nil {
		a.logger.Error("Database connection error", zap.Error(err))
		os.Exit(1)
	}

	a.dbClient = db
	return a
}

func (a *appBuilder) SetCache() *appBuilder {
	db, err := cache.NewCache(a.config.Cache, a.logger).Connect(a.ctx)
	if err != nil {
		a.logger.Error("Cache connection error", zap.Error(err))
		os.Exit(1)
	}

	a.cacheClient = db
	return a
}

func (a *appBuilder) SetServices() *appBuilder {
	a.healthService = services.NewHealthService(a.config, a.logger)

	return a
}

func (a *appBuilder) SetHandler() *appBuilder {
	routes, handlerObj := handler.NewHandler()

	handlerObj.SetHealthHandler(a.healthService)
	a.handler = routes

	return a
}

func (a *appBuilder) Build() (ports.Logger, *fasthttp.Server, int) {
	a.server.Handler = a.handler.Handler
	return a.logger, a.server, a.config.App.Server.Port
}
