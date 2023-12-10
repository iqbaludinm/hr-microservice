package config

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/iqbaludinm/hr-microservice/profile-service/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	host               = utils.GetEnv("DB_HOST")
	port               = utils.GetEnv("DB_PORT")
	username           = utils.GetEnv("DB_USERNAME")
	password           = utils.GetEnv("DB_PASSWORD")
	dbName             = utils.GetEnv("DB_NAME")
	minConns           = utils.GetEnv("DB_MIN_CONNS")
	maxConns           = utils.GetEnv("DB_MAX_CONNS")
	TimeOutDuration, _ = strconv.Atoi(utils.GetEnv("DB_CONNECTION_TIMEOUT"))
)

func NewPostgresDatabase() *pgxpool.Pool {
	logger := utils.NewLogger()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, dbName)

	poolConfig, err := pgxpool.ParseConfig(dsn) // convert dsn
	if err != nil {
		logger.Errorw("Failed to parse configuration", "dsn", dsn)
	}

	minConnsInt, err := strconv.Atoi(minConns)
	if err != nil {
		logger.Errorf("DB_MIN_CONNS expected to be integer", "minimum connections", minConns)
	}
	maxConnsInt, err := strconv.Atoi(maxConns)
	if err != nil {
		logger.Errorf("DB_MAX_CONNS expected to be integer", "maximum connections", maxConns)
	}

	poolConfig.MinConns = int32(minConnsInt)
	poolConfig.MaxConns = int32(maxConnsInt)
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Errorw("Failed to apply pool configuration", "dsn", dsn)
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := pool.Ping(c); err != nil {
		logger.Error(err)
	}

	logger.Infow("Database connected", "dsn", dsn)

	return pool
}
