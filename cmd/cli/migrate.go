package main

func doMigrate(arg2, arg3 string) error {
	// dsn := getDSN()
	checkForDB()

	tx, err := vel.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	// Run the migration command
	switch arg2 {
	case "up":
		// err := vel.MigrateUp(dsn)
		err := vel.RunPopMigrations(tx)
		if err != nil {
			return err
		}
	case "down":
		if arg3 == "all" {
			err := vel.PopMigrateDown(tx, -1)
			if err != nil {
				return err
			}
		} else {
			err := vel.PopMigrateDown(tx, 1)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := vel.PopMigrateReset(tx)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}
	return nil
}
