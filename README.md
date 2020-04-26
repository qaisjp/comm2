# comm2

Have at it. Most of this code is shit. I want to rewrite most of it. Aaaaaaaaaaaaaaa. comm2 will never be released.

## Dependencies

Only Linux and macOS are supported as development environments.

- [Go](https://golang.org/)
    - macOS users: `brew install golang`.
    - Other users: [Getting Started - The Go Programming Language](https://golang.org/doc/install.html).
- [PostgreSQL](https://www.postgresql.org/)
    - macOS users: `brew install postgresql`.
    - Other users: good luck.
- [Node.js](https://nodejs.org/)
    - macOS users: `brew install node`.
    - Other users: [Installing Node.js via package manager | Node.js](https://nodejs.org/en/download/package-manager/).
    - Note that Angular requires a current, active LTS, or maintenance LTS version of Node.js.
- [Angular CLI](https://angular.io/cli)
    - Read [Angular - Setting up the local environment and workspace](https://angular.io/guide/setup-local).
    - Everyone: `npm install -g @angular/cli`

**PostgreSQL on macOS**

If you installed PostgreSQL using Homebrew, you probably should read `brew services --help`.

- Run until manual stop or logout: `brew services run postgresql`
- Run until manual stop + autostart on login: `brew services start postgresql`
- Stop: `brew services stop postgresql`

After upgrading Homebrew's PostgreSQL, you can run `brew postgresql-upgrade-database` to actually upgrade your database.

## Development

1. First clone this repository anywhere:
    ```bash
    git clone git@github.com:multitheftauto/comm2.git
    ```
1. Change to the repository directory: `cd comm2`.

### Database

**Initial setup**

1. Run `psql -U "$(whoami)" postgres` to enter the PostgreSQL interactive terminal, then run these commands:
    ```sql
    -- Create a role/user named `mta`.
    create role mta with login;

    -- Create database `mtahub_dev`
    create database mtahub_dev owner mta;

    -- Make sure database has correct timezone
    alter database mtahub_dev set timezone to 'UTC';
    ```
1. Exit the PostgreSQL interactive terminal.
1. Run `make reset_schema` to import `schema.sql` into database.
1. Run `make migrate` and confirm it says the following:

    ```
    ➜ make migrate
    migrate -path database/migrations -database "postgres://mta@localhost:5432/mtahub_dev?sslmode=disable" up
    no change
    make schema.sql
    pg_dump --no-owner --schema-only -U mta mtahub_dev > schema.sql
    pg_dump --no-owner --data-only -t schema_migrations -U mta mtahub_dev >> schema.sql
    truncate -s -1 schema.sql
    Schema has been written to file 'schema.sql'
    ```

    Running `git diff` now should show NO CHANGES to `schema.sql`,
    otherwise we have forgotten to update `schema.sql`,
    or you have modified the database directly.

**Regenerate schema**

The `schema.sql` contains the latest database schema, with all migrations pre-applied.

Run `make schema.sql` to update this file.

**Migrations**

- `make migrate`: applies migrations from the `database/migrations` folder, then updates `schema.sql`. Useful after a `git pull` that changes the database schema.
- `make NAME=add_reviews migrate_new`: creates a new migration file with the given name.

TODO: Fix poor developer experience when working on migrations asynchronously.

**Other tools**

- `make reset_schema` drops the `mtahub_dev` database and recreate it using the `schema.sql` file.

 It assumes the database is stored locally and you have admin access.

 TODO: the command here doesn't line up with the initial setup. It forgets the timezone change. It's worth just replacing the initial setup with a single call to `make reset_schema`.
- `make checkpoint` _fully_ dumps the database to the `dev_backup` folder.

    Note that this is a full dump, which includes the schema. It is not a data-only dump. The dump is not transferrable between PostgreSQL versions.

    TODO: [add support for data-only dumps](https://github.com/teamxiv/growbot-api/blob/master/Makefile#L35-L42)

- `make restore_checkpoint` restores the most recent _full_ dump in the `dev_backup` folder.

### Backend

**Initial setup**

1. Copy `config.yaml.example` to `config.yaml`. Defaults should be fine.
1. Run `go build ./cmd/mtahub-api` to check that API builds fine.

**Run the server**

1. Start the API with `config=config.yaml ./mtahub-api`.

**Tools**

Auto-refresh is achieved the [gin command line utility](https://github.com/codegangsta/gin), not to be confused with the [Gin HTTP web framework](https://github.com/gin-gonic/gin/), which we also use.

TODO: move to a better task runner.

1. Install `gin` using:
    ```bash
    go get github.com/codegangsta/gin
    ```
1. Verify `gin` was installed correctly by running: `gin -h`.

    If "gin could not be found", you probably need to add `$GOPATH/bin` to your `$PATH`.

    Note that the default `$GOPATH` is `~/go`. You can check that it built fine by looking inside `~/go/bin` for the `gin` binary.
1. Run `make dev_run` and the API should automatically recompile and run when your code changes.

### Frontend

1. Change to the `website` directory: `cd website`.
1. Run `npm install` to install dependencies.
1. Run `ng build` to check that the website builds fine.

**Tools**

Students can use [WebStorm](https://www.jetbrains.com/webstorm/) for free. We have no recommendations.

- Run `ng serve` to "build and serve your app, rebuilding on file changes".
- Use `ng generate` to generate components, services and modules.
- See [Angular CLI](https://angular.io/cli) for more commands.

## License

This repository does not have a license. This means that you don't have the right to do anything with it. This _will_ be open source in the future, but not right now.
