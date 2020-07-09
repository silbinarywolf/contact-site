# Contact Site


# Local Development

1) Start the SQL server only
```
docker-compose up -d db && docker-compose logs -f
```

# Testing

```
GOOS=linux go build -o server && docker-compose --verbose build --no-cache && docker-compose stop && docker-compose up -d && docker-compose logs -f
```

# Setup Developer Environment on Windows with Docker Toolbox

1) Create a virtual machine
```
docker-machine create default
```

2) Setup your environment in your console to use the newly created VirtualMachine
```
eval $(docker-machine env default)
```

3) Build Go binary
```
GOOS=linux go build -o server
```

3) Build and Run environment
```
docker-compose --verbose build --no-cache app && docker-compose stop && docker-compose up -d && docker-compose logs -f
```

**Side Notes*

* Use `docker-machine ip` to get the machines IP address

## Dependencies Considered

- https://github.com/lib/pq
	- On the surface this library seems reasonable but the lack of support when it comes to resolving issues or accepting pull-requests means that we'd have to probably fork this if we hit problems with it, rather than being able to contribute fixes back upstream.
