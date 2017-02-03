# START WITH THE NODE CONTAINER
FROM node:7

MAINTAINER Ramon Saccilotto <ramon.saccilotto@usb.ch>
LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2015"

# ADD THE NODE-MODULES FROM THE BUILD DIRECTORY
COPY node_modules /usr/local/lib/node_modules

# COPY THE HOT RELOAD UTILITY INTO THE BIN DIRECTORY
COPY hot-reload_linux_amd64 /bin/hot-reload

# THIS WILL SYMLINK THE GLOBAL NODE-MODULES INTO OUR WEB DIRECTORY
# ln -f -s /usr/local/lib/node_modules/* /app/web/node_modules/

# MOUNT THE APPLICATION IN THE APP DIRECTORY
VOLUME ["/app"]

# EXPOSE THE PORT 8080 FOR EXTERNAL CONNECTIONS
EXPOSE 8080

# SYMLINK THE GLOBAL NODE-MODULES INTO THE WEBPACK NODE-MODULES DIRECTORY
CMD ["/bin/hot-reload"]
