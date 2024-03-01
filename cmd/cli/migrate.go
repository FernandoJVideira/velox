package main

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	// Run the migration command
	switch arg2 {
	case "up":
		err := vel.MigrateUp(dsn)
		if err != nil {
			return err
		}
	case "down":
		if arg3 == "all" {
			err := vel.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := vel.Steps(dsn, -1)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := vel.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = vel.MigrateUp(dsn)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}
	return nil
}
