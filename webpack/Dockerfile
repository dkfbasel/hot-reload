# START WITH THE NODE CONTAINER
FROM node:8

MAINTAINER Ramon Saccilotto <ramon.saccilotto@usb.ch>
LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2015"

# ADD THE NODE-MODULES FROM THE BUILD DIRECTORY TO THE GLOBAL YARN DIRECTORY
COPY node_modules /usr/local/lib/node_modules

# NOTE: we might later switch to yarn using the yarn global node module path
# COPY node_modules /usr/local/share/.config/yarn/global/node_modules

# COPY THE HOT RELOAD UTILITY INTO THE BIN DIRECTORY
COPY hot-reload/hot-reload_linux_amd64 /bin/hot-reload

# MOUNT THE APPLICATION IN THE APP DIRECTORY
VOLUME ["/app"]

# EXPOSE THE PORT 8080 FOR EXTERNAL CONNECTIONS
EXPOSE 8080

# SYMLINK THE GLOBAL NODE-MODULES INTO THE WEBPACK NODE-MODULES DIRECTORY
CMD ["/bin/hot-reload"]
