package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func doAuth() error {
	//migrations
	dbType := vel.DB.DbType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := vel.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := vel.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table if exists users cascade; drop table if exists tokens cascade;drop table if exists remember_tokens;"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	//run migrations
	err = doMigrate("up", "")
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
