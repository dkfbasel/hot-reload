# Description of npm packages used for development

## autoprefixer
Parse css and add vendor prefixes to rules by Can I Use, This requires the
postcss-loader to work correctly.
https://github.com/postcss/autoprefixer

## babel-core
Javascript transpiler for ES2015.
https://babeljs.io/docs/core-packages/

## babel-loader
Allow transpiling of javascript es2015 files using Babel with webpack.
https://github.com/babel/babel-loader

## babel-plugin-transform-runtime --> remove
Combine all babel helpers from several files. Should be used with babel-runtime.
https://github.com/babel/babel/tree/master/packages/babel-plugin-transform-runtime

## babel-polyfill
Emulate a full ES2015+ environment for older browsers
https://babeljs.io/docs/usage/polyfill/

## babel-preset-minify
Minify javascript that was transpiled with babel
https://github.com/babel/minify

## babel-preset-env
A Babel preset that compiles ES2015+ down to ES5 by automatically determining the Babel
plugins and polyfills you need based on your targeted browser or runtime environments
https://github.com/babel/babel/tree/master/packages/babel-preset-env

## babel-preset-stage-0
Latest version of babel presets with changes to the javascript language that
haven't been approved to be part of a release of Javascript

## babel-runtime --> remove
Externalise references to helpers and builtins, automatically polyfilling your
code without polluting globals. Use with babel-plugin-transform-runtime.
https://babeljs.io/docs/plugins/transform-runtime/

## buble (currently not included)
Fast es2015 compiler. Use as replacement for babel.
https://www.npmjs.com/package/vue-template-loader

## buble-loader (currently not included)
Use buble with webpack.
https://github.com/sairion/buble-loader

## css-nano
Minify css files.
http://cssnano.co

## css-loader
Load css files from javascript. Required to bundle css in webpack
https://github.com/webpack-contrib/css-loader

## cross-env
Set environment variables cross platform. Mainly used to set NODE_ENV to production.
https://github.com/kentcdodds/cross-env

## extract-text-webpack-plugin
Extract text from bundles into a file. Mainly used to extract css from single
file components or separate files into one css file.
https://github.com/webpack-contrib/extract-text-webpack-plugin

## file-loader
Enable hashes for file-paths to improve caching.
https://github.com/webpack-contrib/file-loader

## postcss-loader
Postcss loader for webpack. Postcss is required to enable autoprefixer.
https://github.com/postcss/postcss-loader

## poststylus
Enable autoprefixer with stylus preprocessor.
https://github.com/seaneking/poststylus

## style-loader
Add css to the DOM by injecting a style tag. Required to write styles in single
file components and extract styles with extract-text-webpack-plugin.
https://github.com/webpack-contrib/style-loader

## stylus
Pre-processor for css.
http://stylus-lang.com

## stylus-loader
Load stylus files with webpack. Stylus package needs to be imported as well
https://github.com/shama/stylus-loader

## template-html-loader
Loading html files with templating languages.
https://github.com/jtangelder/template-html-loader

Alternative would be to use html-loader to include hashed paths.
https://github.com/webpack-contrib/html-loader

## vue-loader
Load and compile vuejs files. Includes vue-style-loader.
https://vue-loader.vuejs.org/en/

## vue-svg-loader
Load svg files directoy as vue components. Awesome for icons.
https://github.com/visualfanatic/vue-svg-loader

## vue-template-compiler
Required to precompile vue templates.
https://www.npmjs.com/package/vue-template-loader

## uglifyjs-webpack-plugin
Minification of javscript files
https://github.com/webpack-contrib/uglifyjs-webpack-plugin

## webpack
Webpack utility to transpile the files
https://webpack.js.org

## webpack-cli
Command line interface for webpack
https://webpack.js.org/api/cli/

## webpack-dev-server
Development server with file watching and hot reload
https://webpack.js.org/guides/development/
