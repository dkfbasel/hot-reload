# START WITH THE NODE CONTAINER
FROM node:6

MAINTAINER Ramon Saccilotto <ramon.saccilotto@usb.ch>
LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2015"

# NOTE: DIFFERENT RUN COMMANDS WILL BE CACHED INDEPENDENTLY

# INSTALL WEBPACK, THE WEBPACK-DEV-SERVER AND CROSS-ENV
RUN npm install -g \
    webpack@1.13.1

RUN npm install -g \
    webpack-dev-server@1.14.1

RUN npm install -g \
    cross-env@2.0.0

# INSTALL BABEL FOR COMPILE JAVASCRIPT
RUN npm install -g \
    babel-loader@6.2.4 \
    babel-core@6.13.2 \
    babel-plugin-transform-runtime@6.12.0 \
    babel-preset-es2015@6.13.2 \
    babel-preset-stage-0@6.5.0 \
    babel-runtime@6.11.6

# INSTALL SOME WEBPACK LOADERS
RUN npm install -g \
    file-loader@0.9.0 \
    style-loader@0.13.1 \
    css-loader@0.23.1 \
    stylus-loader@2.1.2 \
    stylus@0.54.5 \
    template-html-loader@0.0.3

# INSTALL THE JEET STYLUS LIBRARY
RUN npm install -g \
    jeet@6.1.4

# INSTALL VUE LOADERS
RUN npm install -g \
    vue-loader@8.5.3 \
    vue-html-loader@1.2.3 \
    vue-style-loader@1.0.0 \
    vue-hot-reload-api@1.3.2

# MOUNT THE APPLICATION IN THE APP DIRECTORY
VOLUME ["/app"]

# TODO: use symlinking for the global node modules
# EXPORT THE GLOBAL NODE PATH
# "/usr/local/lib/node_modules"
# SYMLINK THE GLOBAL NODE-MODULES INTO OUR WEB DIRECTORY
# NOTE: this must be run after the container is linked
# RUN ln -f -s /usr/local/lib/node_modules/* /app/web/node_modules/

# EXPOSE THE PORT 8080 FOR EXTERNAL CONNECTIONS
EXPOSE 8080

# SET THE WORKING DIRECTORY
WORKDIR /app/web
