Hot-reload development for Go and Vue.js
========================================

Note: Please be aware, that the statements and utilities provided are very
opinionated and that you should have some basic understanding on how to use webpack.

Our current development pipeline usually involves a number of different go services
and a single page application front-end written with vue.js. Using docker for
deployment we are able to setup and test projects very quickly and manage a rather large number of services in production with a small number of developers.

Spoiled by the simplicity of using docker for production, we aimed to create a
similar experience for the development. Thus we created two docker containers that can be used to develop such applications making use of hot reloading with minimal setup.

In addition we decided to keep most development dependencies for our front-end
applications contained in the docker container to avoid the necessity of duplicating
a large number of node modules into every project. This also helps with keeping
all modules up to date.


Prerequisites (tested):
-----------------------
- OS X
- Docker for Mac

For building the hot-reload utilities
- Golang installation (with $GOPATH and $GOROOT set)
- gox build utility (only when building the hot-reload utilities)


Usage
-----------------
We generally use the following application structure for our projects and you can
find a scratch project setup using the hot-reload-containers at
https://github.com/dkfbasel/scratch

```
- build         // contains all information required to run the project in production
    - bin       // contains all binaries to run the service
    - public    // contains all files that should be served by the service
    - templates // contains all templates required for the service

- docs          // documentation and asset source files for the project

- src           // contains all development information
    - backend   // contains the golang code for the web server
    - frontend  // contains the web application front-end source code
    - ..        // additional directories for other go packages
    - docker-compose.yml    // configuration for development containers

- test          // directory for separate test data and integration tests

- readme.md     // readme file for every project
```

Docker Compose is used to startup the development services and holds all
configuration required to start hot-reloading for the front and backend.

```
docker-compose up
```

```
docker-compose.yml
------------------

version: '3'

services:

    api:
        image: dkfbasel/hot-reload-go:1.8.1
        ports:
            - "3001:80"
            # goconvey is exposed on port 8080
            - "3002:8080"
        volumes:
            # mount the project into the docker container. Please note, that the
            # app directory is symlinked into the project path specified as
            # environment variable. For goconvey to work, the package must be
            # linked directly into the the package directory i.e. /go/src/[PROJECTPATH]
            - ..:/app
        environment:
            # project is required to make sure that the import paths to
            # optional other packages in the same directory will work as expected
            - PROJECT=github.com/dkfbasel/hot-reload/sample
            # directory is required to set the current directory that should be
            # used for building, it can be omitted if it is the same directory
            # as the project
            - DIRECTORY=src/server
            # specify the command that should be run, can be 'build', 'test' or
            # 'goconvey'. please note, that for goconvey to work the package
            # directory should not be mounted in /app but in /go/src/[PROJECTPATH]
            - CMD=build
            # arguments can be used to specify arguments to pass to the executable
            # on running
            - ARGUMENTS=-test=someString
            # ignore will indicate which directories to ignore from watching
            - IGNORE=/src/web
            
            
    frontend:
        image: dkfbasel/hot-reload-webpack:3.1.0
        # note that the host port and the port on webpack should
        # match to avoid cross origin request issues
        ports:
            - "3000:3000"
        volumes:
            - ..:/app
        environment:
            # directory will be used to define the folder where webpack should
            # be started from and where the local node_modules are to be found
            - DIRECTORY=src/web
            # command defines the command to run after symlinking the global
            # node modules into the local directoy
            - COMMAND=npm run dev

```

It will load the development images specified, which allows for versioning
of the containers, starts up all containers and start the hot-reload development.

Please note that we usually employ additional services in other docker containers,
such as a db container, and connect all services through a network.

All external golang packages should be vendored in the vendor directory. The
container will not attempt to install any go packages. This will ensure full
control of the versioning for the developer and avoid the necessity of having
a working internet connection.

Webpack should be configured to pass api calls forward to the golang backend
service. In addition the host of the webpack dev server should be set to 0.0.0.0
to allow for docker to forward this port to the host (see package.json).

TODO:
----
- [x] Enable hot reloading of go code
- [x] Enable development mode for webpack
- [x] Enable build mode for webpack
- [x] Start services with a simple docker-compose up command

- [x] Enable testing mode for go code
- [x] Support goconvey for testing go code

- [ ] Roadmap: Improve goconvey support (or use a custom test web ui)
- [ ] Roadmap: Allow connection to a cluster of services on a test server
(to avoid the need of starting other micro and db services on the local machine)


Building containers
-------------------

Go: The golang development container will symlink the go package into the directory
specified and watch for changes. Newly added directories will be added to the
watchlist and deleted directories will be removed.

To use it, you should create a corresponding docker container using the
following commands

```
> cd golang

> gox -osarch="linux/amd64" -output="hot-reload_linux_amd64" github.com/dkfbasel/hot-reload/golang/hot-reload

> docker build -t dkfbasel/hot-reload-go:1.8.0 .

> docker run --rm -ti -p 3001:80 -v "$PWD/../sample:/app" -e "PROJECT=github.com/dkfbasel/hot-reload/sample" -e "DIRECTORY=src/server" dkfbasel/hot-reload-go:1.8.0
```

Webpack: The webpack development container will install the node modules specified in
the Dockerfile in the global node directory and symlink all modules into the
local node_modules directory of your project. This is required, as it does
currently not seem to be possible to run webpack from the global directory.

To build the webpack development container make sure the webpack/Dockerfile contains
all node modules you wish to use for your project and follow the steps bellow.

The node modules are copied from the webpack node_modules directory. This allows
to easily check for outdated packages with yarn. Note that it is recommended,
to only include the devDependencies in the webpack container.

```
> cd webpack

> gox -osarch="linux/amd64" -output="hot-reload_linux_amd64" github.com/dkfbasel/hot-reload/webpack/hot-reload

> docker build -t dkfbasel/hot-reload-webpack:3.1.0 .

> docker run --rm -ti -p 3000:3000 -v "$PWD/../sample:/app" -e "DIRECTORY=src/web" -e "COMMAND=npm run dev" dkfbasel/hot-reload-webpack:2.2.0

```
