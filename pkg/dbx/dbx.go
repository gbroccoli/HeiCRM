package dbx

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	dbconf "github.com/gbroccoli/HeiCRM/pkg/config"
	_ "github.com/lib/pq"
)

type configDB struct {
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

var DB *sql.DB

func DSN() string {
	config := &configDB{
		Host:    dbconf.G().Database.Host,
		Port:    dbconf.G().Database.Port,
		User:    dbconf.G().Database.User,
		Pass:    dbconf.G().Database.Pass,
		Name:    dbconf.G().Database.Name,
		SSLMode: dbconf.G().Database.SSLMode,
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(config.User), url.QueryEscape(config.Pass), config.Host, config.Port, config.Name, config.SSLMode,
	)
}

func Open() {
	db, err := sql.Open("postgres", DSN())

	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(45 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err != nil {
		log.Fatal(err)
	}

	DB = db
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func G() *sql.DB {
	if DB == nil {
		log.Panic("DB is nil")
	}
	return DB
}
