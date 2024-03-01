package main

import (
	"fmt"
	"time"
)

func doSessionTable() error {
	// Verify database type
	dbType := vel.DB.DbType

	if dbType == "mariaDB" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}

	fileName := fmt.Sprintf("%d_create_sessions_table", time.Now().UnixMicro())
	upFile := vel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := vel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/"+dbType+"_session.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}
	err = copyDataToFile([]byte("DROP TABLE IF EXISTS sessions;"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	return nil
}
