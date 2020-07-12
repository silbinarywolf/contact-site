# Installation

This document contains instructions for:
- Building the application and necessary artefacts
- Running the application in a Docker environment 

## Install and Setup Docker

1) Create a machine
```
docker-machine create default
```

2) Setup your environment in your console to use the newly created machine
```
eval $(docker-machine env default)
```

If you get the following error and you're using Docker Toolbox, you may need to start up your virtual machine again.
```
Error checking TLS connection: Host is not running
```
