package db

import (
	"database/sql"
	"fmt"
)

type SqlConnectionConfig struct {
	driver           string
	connectionString string
}

func NewSqlConnectionConfig(driver, conn_str string) *SqlConnectionConfig {
	return &SqlConnectionConfig{
		driver:           driver,
		connectionString: conn_str,
	}
}

// func (s *SqlConnectionConfig) InitializeConnection(errChan chan<- error) *sql.DB {
// 	conn, err := sql.Open(s.driver, s.connectionString)
// 	if err != nil {
// 		errChan <- fmt.Errorf("[ERROR]: Sql Connection Could Not Initialized with %v", err)
// 		return nil
// 	}

// 	if err := conn.Ping(); err != nil {
// 		errChan <- fmt.Errorf("[ERROR] DB Ping with %w", err)
// 		return nil
// 	}

// 	return conn
// }

func (s *SqlConnectionConfig) InitializeConnection() (*sql.DB, error) {
	conn, err := sql.Open(s.driver, s.connectionString)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] DB Connection Open Failed with %v", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("[ERROR] DB Connection Ping Failed with %v", err)
	}

	return conn, nil

}
