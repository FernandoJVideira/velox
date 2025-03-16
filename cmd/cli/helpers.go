package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup(arg1, arg2 string) {
	if arg1 != "new" && arg1 != "help" && arg1 != "version" {
		err := godotenv.Load()
		if err != nil {
			exitGracefully(err)
		}

		path, err := os.Getwd()
		if err != nil {
			exitGracefully(err)
		}

		vel.RootPath = path
		vel.DB.DbType = os.Getenv("DATABASE_TYPE")
	}
}

func getDSN() string {
	dbType := vel.DB.DbType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSLMODE"),
			)
		}
		return dsn
	} else {
		return "mysql://" + vel.BuildDSN()
	}
}

func checkForDB() {
	dbType := vel.DB.DbType

	if dbType == "" {
		exitGracefully(errors.New("no database connection in .env file"))
	}

	if !fileExists(vel.RootPath + "/config/database.yml") {
		exitGracefully(errors.New("config/database.yml does not exist"))
	}
}

func showHelp() {
	color.Yellow(`
Velox is a laravel like CLI tool to build web applications.

Usage: 
	velox <command> [arguments]


The commands are:

	help                           - Shows this help message
	new <appname>                  - Creates a new Velox application
	down                           - Puts the Server in maintenance mode
	up                             - Takes the Server out of maintenance mode
	version                        - Shows the current version of the CLI
	migrate                        - Runs all up migrations that have not been run yet
	migrate down                   - Reverts the most recent migration
	migrate reset                  - Runs all down migrations in reverse order and then runs all up migrations
	make migration <name> <format> - Creates up and down migration files in the migrations folder; format = sql/fizz (default fizz)
	make auth                      - Creates the auth tables, files and middleware
	make handler <name>            - Creates a stub handler file in the handlers directory
	make model <name>              - Creates a stub model file in the data directory
	make session                   - Creates a table in the database to store session data
	make key                       - Generates a random 32 character key
	make mail <name>               - Generates 2 starter mail templates in the mail directory
`)
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	//Check for errors
	if err != nil {
		return err
	}
	//Check if the file is a directory
	if fi.IsDir() {
		return nil
	}

	//Check if the file is a .go file
	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}
	//If the file is a .go file, update the file
	if matched {
		//Read the file
		read, err := os.ReadFile(path)
		if err != nil {
			exitGracefully(err)
		}

		newContents := strings.Replace(string(read), "myapp", appURL, -1)

		//Write the file
		err = os.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}

func updateSource() {
	//Walk entire project directory, including subdirectories
	err := filepath.Walk(".", updateSourceFiles)
	if err != nil {
		exitGracefully(err)
	}
}
