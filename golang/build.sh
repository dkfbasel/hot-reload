# !/bin/sh

# build the hot-reload utility
cd ./hot-reload
GOOS=linux GOARCH=amd64 go build -o ../hot-reload_linux_amd64

# build the docker container
cd ../
docker build -t dkfbasel/hot-reload-go:1.13 .