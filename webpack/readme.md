The linker service allows you to run webpack with globally installed packages
inside a docker container.

To use it, you should create a corresponding docker container using the
following commands

```
> cd webpack

> gox -osarch="linux/amd64" -output="linker_linux_amd64" bitbucket.org/dkfbasel/development/webpack/linker

> docker build -t development/webpack .

> docker run --rm -ti -p 8080:8080 -v "$PWD/../_test:/app" -e "DIRECTORY=web" development/webpack

```
