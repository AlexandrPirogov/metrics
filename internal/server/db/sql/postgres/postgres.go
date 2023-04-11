package postgres

import (
	"context"
	"fmt"
	"memtracker/internal/config/server"
	"os"

	"github.com/jackc/pgx/v5"
)

func Ping() error {
	PgURL := server.ServerCfg.DBUrl
	conn, err := pgx.Connect(context.Background(), PgURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}

	defer conn.Close(context.Background())
	err = conn.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}
