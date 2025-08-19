package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bhupendra-dudhwal/go-hexagonal/internal/core/models"
	"github.com/bhupendra-dudhwal/go-hexagonal/internal/core/ports"
	egressPorts "github.com/bhupendra-dudhwal/go-hexagonal/internal/core/ports/egress"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type database struct {
	client *gorm.DB
	logger ports.Logger
	config *models.Database
}

func NewDatabase(config *models.Database, logger ports.Logger) egressPorts.DatabaseConnectionPorts {
	return &database{
		config: config,
		logger: logger,
	}
}

func (d *database) Connect() (*gorm.DB, error) {
	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		d.config.Host,
		d.config.Username,
		d.config.Password,
		d.config.Name,
		d.config.Port,
		d.config.Sslmode,
		d.config.Timezone,
	)

	// Set logger level based on config flag
	var gormLogger logger.Interface
	if d.config.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	var (
		db  *gorm.DB
		err error
	)

	// Retry loop in case DB is not ready yet
	for attempt := 1; attempt <= d.config.ConnectRetries; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
		if err == nil {
			// Verify the connection is actually alive
			sqlDB, _ := db.DB()
			if pingErr := sqlDB.Ping(); pingErr == nil {
				d.logger.Info("Connected to PostgreSQL",
					zap.String("host", d.config.Host),
					zap.Int("port", d.config.Port),
					zap.String("db", d.config.Name),
				)
				break
			} else {
				err = pingErr
			}
		}

		d.logger.Error("DB connection failed ",
			zap.Int("attempt", attempt),
			zap.Int("maxAttempts", d.config.ConnectRetries),
			zap.Error(err),
		)

		if attempt < d.config.ConnectRetries {
			sleep := d.config.RetryInterval
			if sleep <= 0 {
				sleep = time.Second * time.Duration(attempt) // incremental backoff
			}
			time.Sleep(sleep)
		}

	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", d.config.ConnectRetries, err)
	}
	d.client = db

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}
	d.setConnectionPool(sqlDB)

	return d.client, nil
}

func (d *database) setConnectionPool(sqlDB *sql.DB) {
	// Defaults if not provided

	sqlDB.SetMaxIdleConns(d.config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(d.config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(d.config.ConnMaxLife)
	sqlDB.SetConnMaxIdleTime(d.config.ConnMaxIdle)

	d.logger.Info("Connection pool configured:",
		zap.Int("idle", d.config.MaxIdleConns),
		zap.Int("open", d.config.MaxOpenConns),
		zap.Duration("maxLife", d.config.ConnMaxLife),
		zap.Duration("maxLife", d.config.ConnMaxIdle),
	)
}

func (d *database) Close() error {
	if d.client == nil {
		return nil
	}

	sqlDB, err := d.client.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
