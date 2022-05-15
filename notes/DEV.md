# Development guides

All the following commands will assume you are in the project root directory.

The backend server can be run using
```
go run ./backend
```
The server will listen on the address http://localhost:8080.

To build the `chords` CLI, run
```
go build cli/chords.go
```
It can then be invoked using
```
./chords pull [...]
./chords push [...]
./chords backup [...]
```