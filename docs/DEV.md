# Development guide

*NB: all commands shown here will assume you are in the project root directory.*


## Running the server locally

The backend server can be run using
```
go run .
```
There are some environment variables which can be set to affect the behaviour
of the backend:
- `PORT`: the port number which the server will listen on. For example, if
  `PORT=8080`, the server will listen on http://localhost:8080. If `PORT` is
  not set, or set to an invalid value, then the port number will default
  to 8080.
- `LOG_FLAGS`: logging flags which will be sent to the logger. See the
  [log package docs](https://pkg.go.dev/log#pkg-constants) for an explanation.
- `DATABASE_URL`: the address of the database to use (which also encodes the
  type of database).
  - If it's a Postgres URI (`postgres://...`), we'll use the specified Postgres
    database (currently unsupported).
  - If it's empty (`""`), we'll use a temporary database stored in Go memory.
  - Otherwise, we'll treat it as a path on the local filesystem, and use a
    file tree database rooted at that path. See the
    [data model doc](DATA_MODEL.md) for an explanation of the file structure.
- `AUTH_KEY`: the key to use for authorisation on API requests. If the env
  variable is not set, we will try to read the auth key from an `auth_key` file
  in the current working directory.


## Tests

Many of the packages have associated unit tests, which can be run using
`go test ./<package-path>`, e.g.
```
go test ./src/server
```

We also have Go integration tests in the [`tests`](../tests) directory. These
actually start the webserver locally and use the client to connect to the API.
Run these tests using
```
go test ./tests
```

You can run **all** of the unit and integration tests using
```
go test ./...
```


## Deploying to Fly

The production server is hosted on [Fly.io](https://fly.io/) - you can access
it at https://chords.fly.dev/. The `flyctl` command-line tool is used to manage
the deployment. To deploy a new version of the server, run
```
fly deploy
```

Deployment configuration is set in the [fly.toml](../fly.toml) file.


## Command-line interface

I have written a command-line tool that uses the client library to interact
with the API. The source code is located in [src/cmd](../src/cmd). You can
either run it using `go run`:
```
go run ./src/cmd ...
```
or build the binary and run it:
```
go build -o chords ./src/cmd
./chords ...
```

The CLI has methods for each of the API endpoints (to be added), as well as
a `sync` command which copies a local database to a remote server. To update
the production database:
```
./chords sync ./data https://chords.fly.dev/
```