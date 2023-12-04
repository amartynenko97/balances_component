package handler

import (
	"balances_component/database"
	"balances_component/protofile"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

func createAccount(conn *pgxpool.Conn, data *protofile.CreateAccountRequest) error {
	_, err := conn.Exec(context.Background(), database.QueryCreateAccount, data.UserId, data.UserName)
	return err
}
