# Development guides

*NB: all commands shown here will assume you are in the project root directory.*


## Running the backend

The backend server can be run using
```
go run ./backend
```
There are two environment variables which can be set to affect the behaviour of the backend:
- `PORT`: the port number which the server will listen on. For example, if `PORT=8080`, the server will listen on http://localhost:8080. If `PORT` is not set, or set to an invalid value, then the port number will default to 8080.
- `LOG_FLAGS`: logging flags which will be sent to the logger. See the [log package docs](https://pkg.go.dev/log#pkg-constants) for an explanation.


## Running the frontend

You will need to start a live server using a tool such as [this VS Code extension](https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer). After that, navigate to the frontend root directory (`<project-root>/docs/`) in a browser to see the frontend.

The frontend tries to contact the backend at https://barrettj12-chords.herokuapp.com/ (the address of the production server). To use the local server instead, add the following line to your hosts file:
```
127.0.0.1	barrettj12-chords.herokuapp.com
```


## Command-line tool

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