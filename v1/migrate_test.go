package v1

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"net/url"
	"os"
	"testing"

	_ "github.com/ClickHouse/clickhouse-go"
	clickhouseMigrate "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	host     string
	user     string
	password string
	db       string
)

func init() {
	var ok bool

	if host, ok = os.LookupEnv("DB_HOST"); !ok {
		log.Fatal("DB_HOST env is not set")
	}
	if user, ok = os.LookupEnv("DB_USER"); !ok {
		log.Fatal("DB_USER env is not set")
	}
	if password, ok = os.LookupEnv("DB_PASSWORD"); !ok {
		log.Fatal("DB_PASSWORD env is not set")
	}
	if db, ok = os.LookupEnv("DB_NAME"); !ok {
		log.Fatal("DB_NAME env is not set")
	}
}

func TestMigrate(t *testing.T) {
	v := make(url.Values)
	v.Set("username", user)
	v.Set("password", password)
	v.Set("database", db)

	u := url.URL{
		Scheme:   "tcp",
		Host:     host,
		RawQuery: v.Encode(),
	}

	conn, err := sql.Open("clickhouse", u.String())
	if err != nil {
		t.Fatal("failed to connect to DB: ", err)
	}
	defer conn.Close()

	driver, err := clickhouseMigrate.WithInstance(conn, &clickhouseMigrate.Config{})
	if err != nil {
		log.Fatal("failed to create driver: ", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"clickhouse",
		driver,
	)
	if err != nil {
		t.Fatal("failed to create migrate instance: ", err)
	}

	if err := m.Up(); err != nil {
		t.Fatal("failed to migrate: ", err)
	}

	row := conn.QueryRow("SHOW TABLES LIKE 'test_v1'")
	var table string
	if err := row.Scan(&table); err != nil {
		t.Fatal("failed to check table exists: ", err)
	}

	if table != "test_v1" {
		t.Fatal("wrong table name")
	}
}
