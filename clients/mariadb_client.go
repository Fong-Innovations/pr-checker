package clients

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"sync"
	"time"

	goose "github.com/pressly/goose/v3"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type MariaDBClient struct {
	DB *sql.DB
}

// GLOBAL VARIABLE TO HOLD CLIENT
var (
	dbClient *MariaDBClient
	once     sync.Once
)

func NewMariaDBClient(user, password, host string, port int, dbName string) (*MariaDBClient, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to MariaDB: %w", err)
	}
	// Configure connection pool settings
	db.SetConnMaxLifetime(60 * time.Minute) // Maximum lifetime of a connection
	db.SetMaxIdleConns(10)                  // Maximum number of idle connections
	db.SetMaxOpenConns(100)                 // Maximum number of open connections
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MariaDB: %w", err)
	}

	// run db migrations
	goose.SetBaseFS(embedMigrations)
	err = goose.SetDialect("mysql")
	if err != nil {
		log.Println("goose set dialect err:", err)
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		log.Println("goose err:", err)
	}

	log.Println("Successfully connected to MariaDB")
	return &MariaDBClient{DB: db}, nil
}

// NewMariaDBClient creates a new instance of MariaDBClient or returns the existing singleton instance
func GetMariaDBClient(user, password, host string, port int, dbName string) (*MariaDBClient, error) {
	var err error
	once.Do(func() {
		// Check if dbClient is already initialized
		if dbClient != nil {
			log.Println("Creating new MariaDB client")
			dbClient, err = NewMariaDBClient(user, password, host, port, dbName)
		}
	})
	if err != nil {
		log.Println("err creating db client:", err)
		return nil, err
	}

	return dbClient, err
}
