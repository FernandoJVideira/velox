package main

import (
	"github.com/fatih/color"
)

func doAuth() error {
	checkForDB()
	//migrations
	dbType := vel.DB.DbType

	tx, err := vel.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	upBytes, err := templateFS.ReadFile("templates/migrations/auth_tables." + dbType + ".sql")
	if err != nil {
		exitGracefully(err)
	}

	downBytes := ([]byte("drop table if exists users cascade; drop table if exists tokens cascade;drop table if exists remember_tokens;"))
	if err != nil {
		exitGracefully(err)
	}

	err = vel.CreatePopMigration(upBytes, downBytes, "auth", "sql")
	if err != nil {
		exitGracefully(err)
	}

	//run migrations
	err = vel.RunPopMigrations(tx)
	if err != nil {
		exitGracefully(err)
	}
	//copy data
	err = copyFileFromTemplate("templates/data/user.go.txt", vel.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/data/token.go.txt", vel.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/data/remember_token.go.txt", vel.RootPath+"/data/remember_token.go")
	if err != nil {
		exitGracefully(err)
	}

	//Copy middleware
	err = copyFileFromTemplate("templates/middleware/auth.go.txt", vel.RootPath+"/middleware/auth.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt", vel.RootPath+"/middleware/auth-token.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/middleware/remember.go.txt", vel.RootPath+"/middleware/remember.go")
	if err != nil {
		exitGracefully(err)
	}

	//Copy handler
	err = copyFileFromTemplate("templates/handlers/auth-handlers.go.txt", vel.RootPath+"/handlers/auth-handlers.go")
	if err != nil {
		exitGracefully(err)
	}

	//Copy Views
	err = copyFileFromTemplate("templates/mailer/password-reset.html.tmpl", vel.RootPath+"/mail/password-reset.html.tmpl")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/mailer/password-reset.plain.tmpl", vel.RootPath+"/mail/password-reset.plain.tmpl")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/views/login.jet", vel.RootPath+"/views/login.jet")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/views/forgot.jet", vel.RootPath+"/views/forgot.jet")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/views/reset-password.jet", vel.RootPath+"/views/reset-password.jet")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("  - users, tokens, remember_tokens migrations created and executed")
	color.Yellow("  - data/user.go, data/token.go, middleware/auth.go, middleware/auth-token.go created")
	color.Yellow("  - auth middleware created")
	color.Yellow("")
	color.Yellow("Don't forget to add user and token models to your data/models.go and add the appropriate middleware to your routes!")

	return nil
}
