package main

import (
	"database/sql"
	"dukkanim-api/internal/api"
	"dukkanim-api/internal/platform/config"
	"dukkanim-api/internal/platform/db"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	// 1. Configurations

	errChan := make(chan error, 1)

	// Environment Variables Initialization
	cfg := config.Load()

	// 2. PLATFORM INITIALIZATIONS

	// 2.1 DB CONNECTION
	conn, err := initializeDb(cfg)
	if err != nil {
		errChan <- err
		// Os.Exit
	}

	defer conn.Close()

	if err := goose.Up(conn, "migrations/"); err != nil {
		errChan <- err
	}

	server := api.Serve(cfg, conn)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			errChan <- err
		}

	}()
	fmt.Println("Server in alive on: ", cfg.APP_PORT, "🚀......")

	if err := <-errChan; err != nil {
		fmt.Printf("\nSERVER CLOSING DUE TO: \n%v", err)

	}

}

func initializeDb(cfg *config.Config) (*sql.DB, error) {
	// Database Connection Initialization
	db_config := db.NewSqlConnectionConfig("postgres", cfg.CONN_STR)

	return db_config.InitializeConnection()

}
