package clients

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	goose "github.com/pressly/goose/v3"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type MariaDBClient struct {
	DB *sql.DB
}

func NewMariaDBClient(user, password, host string, port int, dbName string) (*MariaDBClient, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, "secret", host, port, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to MariaDB: %w", err)
	}

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
