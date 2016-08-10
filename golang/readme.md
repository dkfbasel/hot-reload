The watchdog service allows you to automatically compile and run a golang package.
It will watch for any file changes and recompile and restart your package.

To use it, you should create a corresponding docker container using the
following commands

```
> cd golang

> gox -osarch="linux/amd64" -output="watchdog_linux_amd64" bitbucket.org/dkfbasel/development/golang/watchdog

> docker build -t development/go .

> docker run --rm -ti -p 8080:80 -v "$PWD/../_test:/app" -e "PROJECT=bitbucket.org/dkfbasel/test" -e "DIRECTORY=server" development/go
```
