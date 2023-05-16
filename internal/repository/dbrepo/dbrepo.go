package dbrepo

import (
	"database/sql"

	"github.com/ASHWIN776/learning-Go/internal/config"
	"github.com/ASHWIN776/learning-Go/internal/repository"
)

// If i want to create a new repo with another database, I can create another struct
type postgresDBRepo struct {
	app *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		app: a,
		DB:  conn,
	}
}
