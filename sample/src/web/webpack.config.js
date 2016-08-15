var path = require('path');
var webpack = require('webpack');

// define some post css plugins to use
var node_modules = path.resolve(__dirname, 'node_modules');
var global_modules = '/usr/local/lib/node_modules';

module.exports = {
	entry: {
		public: [path.resolve(__dirname, 'app/main.js')]
	},
	output: {
		path: path.resolve(__dirname, './public/assets'),
		filename: '[name].bundle.js',
		publicPath: '/assets/'
	},
	devtool: 'source-map',
	devServer: {
		proxy: {
			'/api*': {
				// note that the url to the server is the name
				// of the service that was set in docker-compose.yml
				// it is also possible to use networking and aliases
				target: 'http://api',
				secure: false
			}
		}
	},
	module: {
		loaders: [
			{
				// parse vue components
				test: /\.vue$/,
				loader: 'vue',
				exclude: [node_modules, global_modules]
			}, {
				// edit this for additional asset file types
				test: /\.(png|jpg|gif)$/,
				loader: 'file?name=[name].[ext]?[hash]',
				exclude: [node_modules, global_modules]
			}, {
				// parse css styles
				test: /\.css$/,
				loader: 'style!css!postcss',
				exclude: [node_modules, global_modules]
			}, {
				// parse javascript files
				test: /\.js$/,
				loader: 'babel',
				exclude: [node_modules, global_modules]
			}, {
				// parse stylus styles
				test: /\.styl$/,
				loader: 'style!css!stylus?paths=node_modules/jeet/stylus/',
				exclude: [node_modules, global_modules]
			}
		]
	},
	vue: {
		loaders: {
			stylus: 'style!css!stylus?paths=node_modules/jeet/stylus/',
			exclude: [node_modules, global_modules]
		}
	},
	babel: {
		presets: ['es2015', 'stage-0'],
		plugins: ['transform-runtime']
	}
};
