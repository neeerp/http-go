# HTTP Server

This is an HTTP server I built using the milestones and tests provided on [codecrafters.io](https://codecrafters.io).

You can run the server using the bootstrapping shellscript:
```sh
./your_server.sh --directory <absolute path to serve files from> --port <port to bind to>
```

## Routes
The routes this HTTP server handles are as follows:

- `GET /`: Returns an empty 200
- `GET /echo/:arg`: Echos the argument
- `GET /user-agent`: Echos the User-Agent request header
- `GET /files/:name`: Prints out the contents of the file at the `$directory/$name`
- `POST /files/:name`: Saves the payload as a file at `$directory/$name`

## Design
The design of this server is relatively straight forward.

### Server
The entry point to the server lies in `server.go`. Here is where we bind to the port and fire off go routines to handle
requests in a loop.

### Request Parsing
Request parsing happens in `request.go`. We simply map the raw bytes to a struct that has the status line data, a header
map, and the body.

### Routing and Request Handling
Routing and request handling logic lives in `routes.go`. It's just a switch statement on method + exact or prefix route
URI matches.
