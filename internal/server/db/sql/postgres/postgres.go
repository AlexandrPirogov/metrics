package postgres

import (
	"context"
	"fmt"
	"log"
	"memtracker/internal/config/server"
	"os"

	"github.com/jackc/pgx/v5"
)

func Ping() error {
	PgUrl := server.ServerCfg.DBUrl
	log.Printf("%s", PgUrl)
	conn, err := pgx.Connect(context.Background(), PgUrl)
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
