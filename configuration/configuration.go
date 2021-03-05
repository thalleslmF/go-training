package configuration
import (
	"log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func PrepareDatabase() {
	m, err := migrate.New(
		"file://configuration/db/migrations",
		"postgres://keycloak:password@localhost:5432/keycloak?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}