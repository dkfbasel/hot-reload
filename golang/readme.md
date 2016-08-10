The watchdog service allows you to automatically compile and run a golang package.
It will watch for any file changes and recompile and restart your package.

To use it, you should create a corresponding docker container using the
following commands

> cd golang

> gox -osarch="linux/amd64" -output="watchdog_linux_amd64" bitbucket.org/dkfbasel/development/golang/watchdog

> docker build -t development/go .

> docker run --rm -ti -v "$PWD/../test/server:/app" -e "GOPACKAGE=bitbucket.org/dkfbasel/test" development/go
