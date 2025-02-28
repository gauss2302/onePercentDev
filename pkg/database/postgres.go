package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"math/rand"
	"onepercentdev_server/utils"
	"strings"
	"time"
)

type DB struct {
	*sqlx.DB
}

type DBConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	MaxRetries      int           `mapstructure:"max_retries"`
	RetryDelay      time.Duration `mapstructure:"retry_delay"`
}

var defaultDBConfig = DBConfig{
	MaxOpenConns:    25,
	ConnMaxLifetime: 15 * time.Minute,
	MaxIdleConns:    25,
	ConnMaxIdleTime: 10 * time.Minute,
	MaxRetries:      10,
	RetryDelay:      5 * time.Second,
}

type MigrationFunc func(databaseURL string) error

func NewPostgresDB(databaseURL string, config DBConfig, migrate MigrationFunc) (*DB, error) {
	var db *sqlx.DB
	var err error

	for i := 0; i < config.MaxRetries; i++ {
		db, err = sqlx.Open("postgres", databaseURL)
		if err == nil {
			break
		}

		retryTime := exponentialBackoffWithJitter(i, config.RetryDelay)
		utils.GetLogger().Warn("Failed to connect to database",
			zap.Int("attempt", i+1),
			zap.Int("max_attempts", config.MaxRetries),
			zap.Error(err),
			zap.Duration("retry_time", retryTime),
		)
		time.Sleep(retryTime)
	}

	if err != nil {
		utils.GetLogger().Error("Error connecting to postgres after retries",
			zap.Int("max_retries", config.MaxRetries),
			zap.String("database_url", redactDatabaseURL(databaseURL)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("error connecting to postgres after %d attempts: %w", config.MaxRetries, err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		utils.GetLogger().Error("Error verifying database connection", zap.Error(err))
		return nil, fmt.Errorf("error verifying database connection: %w", err)
	}

	if migrate != nil {
		if err := migrate(databaseURL); err != nil {
			utils.GetLogger().Error("Error running migrations", zap.Error(err))
			return nil, fmt.Errorf("error running migrations: %w", err)
		}
		utils.GetLogger().Info("Migrations ran successfully")
	}

	utils.GetLogger().Info("Connected to postgres")

	return &DB{db}, nil

}

func exponentialBackoffWithJitter(attemp int, baseDelay time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(baseDelay)))
	return baseDelay*time.Duration(1<<uint(attemp)) + jitter
}

func redactDatabaseURL(url string) string {
	redacted := strings.Split(url, "@")
	if len(redacted) > 1 {
		return "postgres://****:****@" + redacted[1]
	}
	return url
}
