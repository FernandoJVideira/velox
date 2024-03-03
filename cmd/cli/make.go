package main

import (
	"errors"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doMake(arg2, arg3, arg4 string) error {
	switch arg2 {
	case "key":
		rnd := vel.RandomString(32)
		color.Yellow("Your new 32 character key is: %s", rnd)
	case "migration":
		checkForDB()

		// dbType := vel.DB.DbType
		if arg3 == "" {
			exitGracefully(errors.New("migration requires a name"))
		}

		//default to migration type of fizz
		migrationType := "fizz"
		var up, down string

		//are we doing fizz or sql?
		if arg4 == "fizz" || arg4 == "" {
			//if fizz read default templates, otherwise use sql
			upBytes, err := templateFS.ReadFile("templates/migrations/migration_up.fizz")
			if err != nil {
				exitGracefully(err)
			}
			downBytes, err := templateFS.ReadFile("templates/migrations/migration_down.fizz")
			if err != nil {
				exitGracefully(err)
			}
			up = string(upBytes)
			down = string(downBytes)
		} else {
			migrationType = "sql"
		}
		//create the migrations for either fizz or sql

		err := vel.CreatePopMigration([]byte(up), []byte(down), arg3, migrationType)
		if err != nil {
			exitGracefully(err)
		}
	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}

	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("you must provide a name for the handler"))
		}
		fileName := vel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}
		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

		err = os.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			exitGracefully(err)
		}

	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("you must provide a name for the model"))
		}

		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		model := string(data)

		plural := pluralize.NewClient()

		var modelName = arg3
		var tableName = arg3

		if plural.IsPlural(arg3) {
			modelName = plural.Singular(arg3)
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(plural.Plural(arg3))
		}

		fileName := vel.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}

		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

		err = copyDataToFile([]byte(model), fileName)
		if err != nil {
			exitGracefully(err)
		}
	case "session":
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}
	case "mail":
		if arg3 == "" {
			exitGracefully(errors.New("you must provide a name for the mail template"))
		}
		htmlMail := vel.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
		plainMail := vel.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"

		err := copyFileFromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
		if err != nil {
			exitGracefully(err)
		}
		err = copyFileFromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}
