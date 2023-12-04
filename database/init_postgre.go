package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

func GetConnection(ctx context.Context, connectionString string, logger *logrus.Logger) (*pgxpool.Pool, error) {
	var dbPool *pgxpool.Pool
	var err error

	for {
		dbPool, err = pgxpool.New(ctx, connectionString)
		if err != nil {
			logger.Errorln("Connection failed! ", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return dbPool, nil
}
