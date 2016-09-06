var path = require('path');
var webpack = require('webpack');

// extract text into separate files
var ExtractTextPlugin = require('extract-text-webpack-plugin');

// define the app directory to include in compilation
var app_directory = path.resolve(__dirname, 'app');

module.exports = {
	entry: {
		public: [path.resolve(__dirname, 'app/main.js')]
	},
	output: {
		path: path.resolve(__dirname, './public/assets'),
		filename: '[name].bundle.js',
		publicPath: '/assets/'
	},
	module: {
		loaders: [
			{
				// parse vue components
				test: /\.vue$/,
				loader: 'vue',
				include: [app_directory]
			}, {
				// edit this for additional asset file types
				test: /\.(png|jpg|gif)$/,
				loader: 'file?name=[name].[ext]?[hash]',
				include: [app_directory]
			}, {
				// parse css styles
				test: /\.css$/,
				loader: 'style!css!postcss',
				include: [app_directory]
			}, {
				// parse javascript files
				test: /\.js$/,
				loader: 'babel',
				include: [app_directory]
			}, {
				// parse stylus styles
				test: /\.styl$/,
				loader: ExtractTextPlugin.extract('style', 'css!stylus?paths=node_modules/jeet/stylus/'),
				include: [app_directory]
			}
		]
	},
	vue: {
		loaders: {
			stylus: ExtractTextPlugin.extract('style', 'css!stylus?paths=node_modules/jeet/stylus/'),
			include: [app_directory]
		}
	},
	babel: {
		presets: ['es2015', 'stage-0'],
		plugins: ['transform-runtime']
	},
	plugins: [
		new ExtractTextPlugin('[name].css'),
		new webpack.DefinePlugin({
			'process.env': {
				NODE_ENV: '"production"'
			}
		}),
		new webpack.optimize.UglifyJsPlugin({
			compress: {
				warnings: false
			}
		}),
		new webpack.optimize.OccurenceOrderPlugin()
	]
};
