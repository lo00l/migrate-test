package v2

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/golang-migrate/migrate/v4"
	clickhouseMigrate "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"net/url"
	"os"
	"testing"
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
	u := url.URL{
		Scheme: "clickhouse",
		Host:   host,
		Path:   db,
		User:   url.UserPassword(user, password),
	}
	s := u.String()
	fmt.Println(s)

	conn, err := sql.Open("clickhouse", u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	driver, err := clickhouseMigrate.WithInstance(conn, &clickhouseMigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"clickhouse",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	row := conn.QueryRow("SHOW TABLES LIKE 'test_v2'")
	var table string
	if err := row.Scan(&table); err != nil {
		t.Fatal("failed to check table exists: ", err)
	}

	if table != "test_v2" {
		t.Fatal("wrong table name")
	}
}
