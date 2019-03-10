# Reflection Server
HTTP request reflection server for understanding requests.

## Usage
Compile the binaries using `make compile`.

You should find the `./bin` directory populated with binaries and their SHA sums. Run the approrpriate one for your platform.

### Configuration

- `--port int`: Defines the port for the HTTP server to run on. The TLS version will run on this port plus one. Example: `./bin/reflection --port 1111` will make the server listen on port 1111 for HTTP and 1112 for HTTPS

- `--version`: Displays the version output

## Development
Run `make tlscerts` to generate the server certificate and key.

Run `godev` to start the live-reload development environment.

Try it out by doing `curl`s to it:

```sh
curl -k \
  --cookies 'A=B,hello=world'
  -X POST \
  -d "hi there!" \
  'https://localhost:8081/some/path?some=variables';
```

# Cheers
