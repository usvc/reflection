# Reflection Server
HTTP request reflection server for understanding and debugging requests.

[![Latest Version](https://badge.fury.io/gh/usvc%2Freflection.svg)](https://github.com/usvc/reflection/releases)

## Usage

### Docker Image
You can also spin this up quickly using docker.

#### Without TLS
Simply run:

```sh
docker run \
  -p 8080:8080 \
  usvc/reflection:latest;
```

#### With TLS
First generate the server certificates:

```sh
make tlscerts
```

The `server.key` and `server.cert` should appear in your working directory.

```sh
docker run \
  -v "$(pwd)/server.key:/server.key" \
  -v "$(pwd)/server.crt:/server.crt" \
  -p 8080:8080 \
  -p 8081:8081 \
  usvc/reflection:latest
```

### Binary
Compile the binaries using `make compile`.

You should find the `./bin` directory populated with binaries and their SHA sums. Run the approrpriate one for your platform:

```sh
# if you want tls
make tlscerts;

# without tls
./bin/reflection-linux-amd64;
```

### Testing

#### Non-TLS
After starting the application, do a `curl` to it and pipe it to `jq` (if you have it):

```sh
curl \
  -vv \
  --cookie 'hello=world' \
  -X POST \
  -d "hi there!" \
  'http://localhost:8080/some/path?some=query&hello=world&hello=everyone' \
  | jq
```

#### TLS
After starting the application, do a `curl` to it allowing for insecure requests and pipe it to `jq` (if you have it):

```sh
curl \
  -vv \
  -k \
  --cookie 'hello=world' \
  -X POST \
  -d "hi there!" \
  'https://localhost:8081/some/path?some=query&hello=world&hello=everyone' \
  | jq
```


### Configuration

- `--port int`: Defines the port for the HTTP server to run on. The TLS version will run on this port plus one. Example: `./bin/reflection --port 1111` will make the server listen on port 1111 for HTTP and 1112 for HTTPS

- `--version`: Displays the version output

## Development
Run `git submodule init` to bring in the Go generators.

Run `make generate` to re-generate the `version.go`

Run `make tlscerts` to generate the server certificate and key.

Run `godev` to start the live-reload development environment.

Try it out by doing `curl`s to it:

```sh
curl -k \
  -vv \
  --cookies 'A=B,hello=world'
  -X POST \
  -d "hi there!" \
  'https://localhost:8081/some/path?some=variables';
```

## Licensing
This project is licensed under the MIT license. See [the LICENSE file](./LICENSE) for the full text.

# Cheers
