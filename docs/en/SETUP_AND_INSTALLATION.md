# Setup and Installation

This document contains instructions for:
- Building the application and necessary artefacts
- Running the application in a Docker environment

# Build Environment

The following assumes you're running a *nix console such as Bash, Git Bash or macOS Terminal.
If you're using Windows, I recommend Git Bash. Don't try to use the Windows Command Line or Powershell, the commands below will not work.

These docs also assume you understand how to setup a remote docker machine and setup your environment so that when the following commands are executed, they'll be done against your remote server.

1) Create a copies of the provided example files

    - [config.example.json](/config.example.json)
    - [docker-compose.prod.example.yml](/docker-compose.prod.example.yml)

Remove the ".example" part from each file and configure them for production use.

2) Update the example files to be more secure and production ready

    - Change the database user and password in both `config.json` and `docker-compose.prod.yml` to not be admin/password.

3) The following command-line statements will:

    - Build Go binary for Linux. (this binary file will be packaged into the Docker container, as defined in Dockerfile)
    - Force build of all images. Mostly used to rebuild the "app" image.
    - Stop the containers if they're running
    - Make the containers run, using our [docker-compose.prod.yml](./docker-compose.prod.yml) override file
    - We run it in detached mode so that if we close the console window, the server will keep running.
```
GOOS=linux go build -o server && 
docker-compose --verbose build --no-cache && 
docker-compose stop && 
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d && 
docker-compose logs -f
```

If you've done things correctly, your logs should end in a line like this:
```
app_1 | 2020/07/12 07:24:58 Starting server on :8080...
```

4) You can use the following command to get the IP address of the Docker machine and visit it in the browser
```
docker-machine ip
```

ie. It might give "192.168.99.100", so you'd visit "http://192.168.99.100:8080" in Chrome.

# Destroying the environment

The following will stop and delete your containers. This means you'll lose all data in your SQL database.

```
docker-compose stop &&
docker-compose rm
```
