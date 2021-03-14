package configuration
import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func PrepareDatabase() *gorm.DB {
	m, err := migrate.New(
		"file://configuration/db/migrations",
		"postgres://keycloak:password@localhost:5432/keycloak?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m.Up())
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	dsn := "host=localhost user=keycloak password=password dbname=keycloak port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}