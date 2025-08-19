package models

import (
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	App      *App      `yaml:"app"`
	Logger   *Logger   `yaml:"logger"`
	Database *Database `yaml:"database"`
	Cache    *Cache    `yaml:"cache"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.App, validation.Required, validation.NotNil),
		validation.Field(&c.Logger, validation.Required, validation.NotNil),
		validation.Field(&c.Database, validation.Required, validation.NotNil),
		validation.Field(&c.Cache, validation.Required, validation.NotNil),
	)
}

type App struct {
	Server *Server `yaml:"server"`
}

func (a App) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Server, validation.Required, validation.NotNil),
	)
}

type Server struct {
	Compression bool                  `yaml:"compression"`
	Environment constants.Environment `yaml:"environment"`
	Port        int                   `yaml:"port"`
}

func (s Server) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Environment, validation.Required),
		validation.Field(&s.Port, validation.Required, validation.Min(1111)),
	)
}

type Logger struct {
	Level string `yaml:"level"`
}

func (l Logger) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Level, validation.Required),
	)
}

type Database struct {
	Debug          bool          `yaml:"debug"`
	Name           string        `yaml:"name"`
	Host           string        `yaml:"host"`
	Username       string        `yaml:"username"`
	Password       string        `yaml:"password"`
	Port           int           `yaml:"port"`
	Timezone       string        `yaml:"timezone"`
	Sslmode        string        `yaml:"sslmode"`
	MaxIdleConns   int           `yaml:"maxIdleConns"`
	MaxOpenConns   int           `yaml:"maxOpenConns"`
	ConnMaxLife    time.Duration `yaml:"connMaxLife"`
	ConnMaxIdle    time.Duration `yaml:"connMaxIdle"`
	ConnectRetries int           `yaml:"connectRetries"`
	RetryInterval  time.Duration `yaml:"retryInterval"`
}

func (d Database) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.Port, validation.Required, validation.Min(1111)),
		validation.Field(&d.Host, validation.Required),
		validation.Field(&d.Username, validation.Required),
		validation.Field(&d.Password, validation.Required),
		validation.Field(&d.Timezone, validation.Required),
		validation.Field(&d.Sslmode, validation.Required),
		validation.Field(&d.MaxIdleConns, validation.Required, validation.Min(10)),
		validation.Field(&d.MaxOpenConns, validation.Required, validation.Min(5)),
		validation.Field(&d.ConnMaxLife, validation.Required),
		validation.Field(&d.ConnMaxIdle, validation.Required),
		validation.Field(&d.ConnectRetries, validation.Required, validation.Min(3)),
		validation.Field(&d.RetryInterval, validation.Required),
	)
}

type Cache struct {
	Name           int           `yaml:"name"`
	Host           string        `yaml:"host"`
	Username       string        `yaml:"username"`
	Password       string        `yaml:"password"`
	Port           int           `yaml:"port"`
	PoolSize       int           `yaml:"poolSize"`
	MinIdleConns   int           `yaml:"minIdleConns"`
	DialTimeout    time.Duration `yaml:"dialTimeout"`
	ReadTimeout    time.Duration `yaml:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout"`
	ConnectRetries int           `yaml:"connectRetries"`
	RetryInterval  time.Duration `yaml:"retryInterval"`
}

func (c Cache) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Port, validation.Required, validation.Min(1111)),
		validation.Field(&c.Host, validation.Required),
		// validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.PoolSize, validation.Required, validation.Min(5)),
		validation.Field(&c.MinIdleConns, validation.Required),
		validation.Field(&c.DialTimeout, validation.Required),
		validation.Field(&c.ReadTimeout, validation.Required),
		validation.Field(&c.WriteTimeout, validation.Required),
		validation.Field(&c.ConnectRetries, validation.Required, validation.Min(3)),
		validation.Field(&c.RetryInterval, validation.Required),
	)
}
