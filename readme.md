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
        image: dkfbasel/hot-reload-go:1.26.1
        ports:
            - "3001:80"
        volumes:
            # Mount the project into the docker container. Must use go modules.
            - ..:/app
            # Use docker volumes to cache builds and go modules. This will speed
            # up subsequent compilation times.
            - gomod:/go/pkg/mod
            - gocache:/root/.cache/go-build
        environment:
            # Directory to look for the main go entry point (default: /app)
            - DIRECTORY=/app
            # Specify the command that should be run, can be 'build' or 'test'
            # 'build' is the default command 
            - CMD=build
            # Arguments can be used to specify arguments to pass to the executable
            # on running
            - ARGS=-test=someString
            # Ignore will indicate which files and subdirectories to ignore from 
            # watching, relative and absolute paths are supported
            - IGNORE=/cmd/web,*.md
            # Watch will add additional files and directories to add to watch
            # for changes. Tthe ignore list is applied to those directories as well
            - WATCH=/src,.env
        
volumes:
  gomod:
  gocache:
```
