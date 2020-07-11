# Contact Site



# Setup Docker

1) Create a virtual machine
```
docker-machine create default
```

2) Setup your environment in your console to use the newly created VirtualMachine
```
eval $(docker-machine env default)
```

If you get the following error and you're using Docker Toolbox, you may need to start up the virtual machine again:
```
Error checking TLS connection: Host is not running
```

# Local Development

1) Start the SQL server only
```
docker-compose up -d db && docker-compose logs -f
```

2) Build and run your binary
```
go build && ./contact-site
```

# Testing

```
GOOS=linux go build -o server && docker-compose --verbose build --no-cache && docker-compose stop && docker-compose up -d && docker-compose logs -f
```

# Run Environment

1) Build Go binary
```
GOOS=linux go build -o server
```

2) Build and Run environment
```
docker-compose --verbose build --no-cache app && docker-compose stop && docker-compose up -d && docker-compose logs -f
```

**Side Notes*

* Use `docker-machine ip` to get the machines IP address

# Run Tests

* Run all tests
```
go test ./...
```

* Run unit tests
```
go test ./internal/...
```

* Run integration tests
```
go test ./test
```

# Run Tests With Code Coverage

* Generate code coverage
```
go test ./test -cover -coverpkg=./... -coverprofile=coverage.out
```

* Open coverage file in the browser
```
go tool cover -html=coverage.out
```

## Dependencies Considered

- https://github.com/lib/pq
	- On the surface this library seems reasonable but the lack of support when it comes to resolving issues or accepting pull-requests means that we'd have to probably fork this if we hit problems with it, rather than being able to contribute fixes back upstream.
