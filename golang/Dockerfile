# START FROM GO ALPINE
FROM golang:1.13-alpine

LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2019"

# ADD ADDITIONAL PACKAGES
# - bash for interactive bash
# - git to download go packages
# - vim to edit files inside the container
RUN apk update && apk upgrade && \
	apk add --no-cache bash git vim

# ADD SSH CLIENT TO CONNECT TO GIT REPOSITORIES VIA SSH AND CUSTOM KEYS
RUN apk add --no-cache openssh-client

# COPY THE HOT RELOAD UTILITY INTO THE BIN DIRECTORY
COPY hot-reload_linux_amd64 /bin/hot-reload

# THE PROJECT TO WATCH SHOULD BE CONNECTED ON THE /APP VOLUME
VOLUME ["/app"]

# EXPOSE PORT 80 FOR EXTERNAL CONNECTIONS
EXPOSE 80

# WATCH FOR CHANGES AND AUTOMATICALLY REBUILD THE APPLICATION
CMD ["/bin/hot-reload"]
