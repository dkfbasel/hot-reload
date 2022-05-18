# --- BUILD IMAGE ---

# USE GOLANG CONTAINER TO BUILD THE APPLICATION
FROM golang:1.18.2 as build

WORKDIR /src

# COPY GO MODULE FILES TO ALLOW FOR CACHING OF MODULE FETCHING
COPY ./go.mod /src/go.mod
COPY ./go.sum /src/go.sum

RUN go mod download

# COPY THE SOURCE CODE INTO THE CONTAINER
COPY . /src

# BUILD THE HOT RELOAD UTILITY (CGO MUST BE DISABLED TO RUN IT ON ALPINE)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ./hot-reload

# --- RUNTIME IMAGE ---

# START FROM GO ALPINE
FROM golang:1.18.2-alpine3.15

LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2022"

# ADD ADDITIONAL PACKAGES
# - bash for interactive bash
# - git to download go packages
# - vim to edit files inside the container
RUN apk update && apk upgrade && \
	apk add --no-cache bash git vim

# ADD CURL
RUN apk add --update curl && \
	rm -rf /var/cache/apk/*

# ADD SSH CLIENT TO CONNECT TO GIT REPOSITORIES VIA SSH AND CUSTOM KEYS
RUN apk add --no-cache openssh-client

# COPY THE HOT RELOAD UTILITY INTO THE BIN DIRECTORY
COPY --from=build /src/hot-reload /bin/hot-reload

# THE PROJECT TO WATCH SHOULD BE CONNECTED ON THE /APP VOLUME
VOLUME ["/app"]

# EXPOSE PORT 80 FOR EXTERNAL CONNECTIONS
EXPOSE 80

# WATCH FOR CHANGES AND AUTOMATICALLY REBUILD THE APPLICATION
CMD ["/bin/hot-reload"]
