# Velox

An open source Laravel alternative built in Go.

## Description

Velox is able to map almost all of Laravel's functionalities such as:

- CLI for easy project creation & management
- Web Page Rendering
- Support for different database types (mySQL/MariaDB and Postgres)
- Database Migration Support (SQl & Soda Migrations)
- Session Management & Multiple Session Storage options (Cookie, Redis, mySQL, or Postgres)
- Cache management (Badger or Redis)
- CSRF Protection
- Emailing System
- Full Auth System (w/SSO Login Support)
- Remote File Systems Support (Minio, sFTP, WebDAV, Amazon S3 Buckets)
- RPC Support
- Graceful Shutdown
- Easy to use testing utilities (Similar to Laravel Dusk)

## Instalation

**Prerequisites:** Having [Make](https://www.gnu.org/software/make/) and [Go (1.22)](https://go.dev/) Installed

Clone this repo or download the contents of the repo:

```
git clone https://github.com/FernandoJVideira/velox.git
```

On Windows:
Open the `Makefile` and change line 15 to:

```
@go build -o dist/velox.exe ./cmd/cli
```

Here you will have to add the dist directory (or any directory that holds the executable) to your Path Environment Variable in order to use the commands.

MacOS/Linux:

Change directories to the cloned/downloaded folder and type the following command:

```
make build
```

The executable by default is located in the /dist directory inside the cloned folder. Either copy it to the desired path (in this case you can run all commands with `./velox` and inside the directory that holds the executable) or add it to the $PATH variable to have access to the `velox` command.

To create a new project you can use:

```
velox new <project_name>
```

To get info about all the commands supported by the CLI:

```
velox help
```

Here's what the `help` command outputs:

```
Available commands:

	help                           - Shows this help message
	new <appname>                  - Creates a new Velox application
	down                           - Put the Server in maintenance mode
	up                             - Take the Server out of maintenance mode
	version                        - Shows the current version of the CLI
	migrate                        - Runs all up migrations thet have not been run yet
	migrate down                   - Reverts the most recent migration
	migrate reset                  - Runs all down migrations in reverse order and then runs all up migrations
	make migration <name> <format> - Creates up and down migration files in the migrations folder; format = sql/fizz (default fizz)
	make auth                      - Creates the auth tables, files and middleware
	make handler <name>            - Creates a stub handler file in the handlers directory
	make model <name>              - Creates a stub model file in the data directory
	make session                   - Creates a table in the database to store session data
	make key                       - Generates a random 32 character key
	make mail <name>               - Generates 2 starter mail templates in the mail directory
```

## Examples and docs (WIP, will be commited soon)
