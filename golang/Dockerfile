# START WITH THE GOLANG ALPINE PACKAGE
FROM golang:1.6-alpine

MAINTAINER Ramon Saccilotto <ramon.saccilotto@usb.ch>
LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2015"

# ADD A BASH SHELL
RUN apk add --update bash && rm -rf /var/cache/apk/*

# COPY THE HOT RELOAD UTILITY INTO THE BIN DIRECTORY
COPY hot-reload_linux_amd64 /bin/hot-reload

# ADD ROOT CERTIFICATE PROVIDERS FOR SSL CONNECTIONS
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# THE PROJECT TO WATCH SHOULD BE CONNECTED ON THE /APP VOLUME
VOLUME ["/app"]

# EXPOSE PORT 80 FOR EXTERNAL CONNECTIONS
EXPOSE 80

# SYMLINK THE APPLICATION DIRECTORY INTO THE GOLANG SOURCE, WATCH FOR CHANGES
# AND AUTOMATICALLY REBUILD THE APPLICATION
CMD ["/bin/hot-reload"]
