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
        image: dkfbasel/hot-reload-go:1.16.2
        ports:
            - "3001:80"
        volumes:
            # mount the project into the docker container. Must use go modules.
            - ..:/app
            # mount modules directory from source code or as docker volume to
            # cache go modules
            - ../_modules:/go/pkg/mod
        environment:
            # directory to look for the main go entry point (default: /app)
            - DIRECTORY=/app
            # specify the command that should be run, can be 'build' or 'test'
            # 'build' is the default command 
            - CMD=build
            # arguments can be used to specify arguments to pass to the executable
            # on running
            - ARGS=-test=someString
            # ignore will indicate which subdirectories to ignore from watching
            - IGNORE=/src/web
        
```
