# START WITH THE GOLANG ALPINE PACKAGE
FROM golang:1.8-alpine

MAINTAINER Ramon Saccilotto <ramon.saccilotto@usb.ch>
LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2015"

# ADD ROOT CERTIFICATE PROVIDERS FOR SSL CONNECTIONS
RUN apk add --no-cache ca-certificates

# ADD GIT TO DOWNLOAD GO DEPENDENCIES
RUN apk update && apk upgrade && \
	apk add --no-cache bash git

# INSTALL GOCONVEY FOR TESTS
RUN go get github.com/smartystreets/goconvey

# COPY THE HOT RELOAD UTILITY INTO THE BIN DIRECTORY
COPY hot-reload_linux_amd64 /bin/hot-reload


# THE PROJECT TO WATCH SHOULD BE CONNECTED ON THE /APP VOLUME
VOLUME ["/app"]

# EXPOSE PORT 80 FOR EXTERNAL CONNECTIONS
EXPOSE 80 8080

# SYMLINK THE APPLICATION DIRECTORY INTO THE GOLANG SOURCE, WATCH FOR CHANGES
# AND AUTOMATICALLY REBUILD THE APPLICATION
CMD ["/bin/hot-reload"]
