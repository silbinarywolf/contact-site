# Contributing

This document contains instructions for:
- Setting up an environment locally for development
- Running test code


## Fast/Iterative Local Development

Changing Go code, rebuilding and then restarting the entire Docker environment can waste a fair bit of time. To speed up iteration times, I just run only the PostgresSQL server 

1) Start the PostgresSQL server in detached mode (-d) and to just immediately show me the logs. I do things this way so that if you close the console window running, your SQL server will continue to run in the background.
```
docker-compose up -d db && docker-compose logs -f
```

2) Build and run your binary
```
go build && ./contact-site
```

## Run Tests

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

## Run Tests With Code Coverage

* Generate code coverage
```
go test ./test -cover -coverpkg=./... -coverprofile=coverage.out
```

* Open coverage file in the browser
```
go tool cover -html=coverage.out
```
