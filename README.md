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
- [Angular CLI](https://angular.io/)
    - Read [Angular - Setting up the local environment and workspace](https://angular.io/guide/setup-local).
    - Everyone: `npm install -g @angular/cli`

**PostgreSQL on macOS**

If you installed PostgreSQL using Homebrew, you probably should read `brew services --help`.

- Run until manual stop or logout: `brew services run postgresql`
- Run until manual stop + autostart on login: `brew services start postgresql`
- Stop: `brew services stop postgresql`

After upgrading Homebrew's PostgreSQL, you can run `brew postgresql-upgrade-database` to actually upgrade your database.

## Development

**Initial setup**

1. Clone this repository anywhere, `git clone git@github.com:multitheftauto/comm2.git`.
2. PostgreSQL stuff
    ```bash
    # Access postgres 'shell'
    psql -U "$(whoami)" postgres

    # Create a role/user named `mta`.
    create role mta with login;

    # Create database `mtahub_dev`
    create database mtahub_dev owner mta;

    # Make sure database has correct timezone
    alter database mtahub_dev set timezone to 'UTC';
    ```
3. Copy `config.yaml.example` to `config.yaml`. Defaults should be fine.
4. Run `make reset_schema` to import schema into database.
5. Run `make migrate` and confirm it says the following:

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

    Running `git diff` should show NO CHANGES to `schema.sql`,
    otherwise we have forgotten to update `schema.sql`,
    or you have modified the database directly.
6. **TODO**: frontend docs

**Run the API**

1. Run `go build ./cmd/mtahub-api` to build the API
2. Start the API with `config=config.yaml ./mtahub-api`

**Run the frontend**

Todo, just do `ng serve` or something.

## OAuth

- Send the user to `https://forum.mtasa.com/oauth/authorize/?client_id={CLIENT_ID}&response_type=code&redirect_uri=http://localhost:8080`

## License

This repository does not have a license. This means that you don't have the right to do anything with it. This _will_ be open source in the future, but not right now.
