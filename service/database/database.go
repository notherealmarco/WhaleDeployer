/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/notherealmarco/WhaleDeployer/service/structures"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// GetProjects returns all the projects
	GetProjects() (*[]structures.Project, error)
	// GetProject returns the project with the given ID
	GetProject(id string) (*structures.Project, error)
	// AddProject adds a new project
	AddProject(link *structures.Project) error
	// DeleteProject deletes the project with the given ID
	DeleteProject(name string) error

	BuildProject(name string) error
	BuildSuccess(name string) error
	BuildFail(name string) error

	Ping() error
}

// DBTransaction is the interface for a generic database transaction
type DBTransaction interface {
	// Commit commits the transaction
	Commit() error
	// Rollback rolls back the transaction
	Rollback() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='projects';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE "projects" (
			"name"	TEXT NOT NULL,
			"path"	TEXT NOT NULL,
			"description"	TEXT,
			"git_url" TEXT,
			"git_branch" TEXT,
			"dockerfile"	TEXT NOT NULL,
			"image_name"	TEXT NOT NULL,
			"image_tag"	TEXT NOT NULL DEFAULT 'latest',
			"deploy_key"	INTEGER NOT NULL DEFAULT 0,
			"last_build"	TEXT NOT NULL DEFAULT '1970-01-01',
			"status"	TEXT NOT NULL DEFAULT 'unknown',
			PRIMARY KEY("name")
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
