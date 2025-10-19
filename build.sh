# !/bin/sh

# build the hot reload utility and a respective docker container
docker build -t dkfbasel/hot-reload-go:1.25.3 -f ./Dockerfile ./cmd/hot-reload
