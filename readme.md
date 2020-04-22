Hot-reload Development for Go in Docker Containers
==================================================

This directory contains the source code for the image dkfbasel/hot-reload-go. It
will compile and start the go program linked into the container specified under 
directory (/app per default) and automatically recompile and reload the program 
when any file changes.

Please note that go modules is required for it to work.

```
docker-compose.yml
------------------

version: '3'

services:

    api:
        image: dkfbasel/hot-reload-go:1.14.2
        ports:
            - "3001:80"
        volumes:
            # mount the project into the docker container. Please note, that the
            # app directory is symlinked into the project path specified as
            # environment variable. For goconvey to work, the package must be
            # linked directly into the the package directory i.e. /go/src/[PROJECTPATH]
            - ..:/app
        environment:
            # directory to look for the main go entry point (default: /app)
            - DIRECTORY=/app
            # specify the command that should be run, can be 'build' or 'test'
            # 'build' is the default command 
            - CMD=build
            # arguments can be used to specify arguments to pass to the executable
            # on running
            - ARGUMENTS=-test=someString
            # ignore will indicate which subdirectories to ignore from watching
            - IGNORE=/src/web
        
```
