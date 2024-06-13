package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgreDatabase(cfg *viper.Viper) (db *gorm.DB, err error) {
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		cfg.GetString("POSTGRE_HOST"),
		cfg.GetString("POSTGRE_USERNAME"),
		cfg.GetString("POSTGRE_PASSWORD"),
		cfg.GetString("POSTGRE_NAME"),
		cfg.GetInt("POSTGRE_PORT"),
	)
	db, err = gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(getLoggerLevel(cfg.GetString("GORM_LOGGER"))),
	})
	if err != nil {
		return
	}

	pgsqlDB, err := db.DB()
	if err != nil {
		return
	}

	err = pgsqlDB.Ping()
	if err != nil {
		return
	}

	pgsqlDB.SetMaxOpenConns(cfg.GetInt("POSTGRE_POOL_SIZE"))
	pgsqlDB.SetConnMaxLifetime(cfg.GetDuration("POSTGRE_MAX_CONN_LIFETIME"))
	pgsqlDB.SetMaxIdleConns(cfg.GetInt("POSTGRE_MAX_IDLE_CONNS"))

	return
}

// getLoggerLevel return gorm log level setup based from env/config of the app
func getLoggerLevel(v string) logger.LogLevel {
	switch v {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Silent
	}
}
