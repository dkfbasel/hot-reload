var path = require('path');
var webpack = require('webpack');

// extract text into separate files
var ExtractTextPlugin = require('extract-text-webpack-plugin');

// define some post css plugins to use
var node_modules = path.resolve(__dirname, 'node_modules');

module.exports = {
	entry: {
		public: [path.resolve(__dirname, 'app/public/main.js')],
		internal: [path.resolve(__dirname, 'app/internal/main.js')],
		admin: [path.resolve(__dirname, 'app/admin/main.js')],
	},
	output: {
		path: path.resolve(__dirname, '../_build/web/assets'),
		filename: '[name].bundle.js',
		publicPath: '/assets/'
	},
	module: {
		loaders: [
			{
				// parse vue components
				test: /\.vue$/,
				loader: 'vue',
				exclude: node_modules
			}, {
				// edit this for additional asset file types
				test: /\.(png|jpg|gif)$/,
				loader: 'file?name=[name].[ext]?[hash]',
				exclude: node_modules
			}, {
				// parse css styles
				test: /\.css$/,
				loader: 'style!css!postcss',
				exclude: node_modules
			}, {
				// parse javascript files
				test: /\.js$/,
				loader: 'babel',
				exclude: node_modules
			}, {
				// parse stylus styles
				test: /\.styl$/,
				loader: ExtractTextPlugin.extract('style', 'css!stylus?paths=node_modules/jeet/stylus/'),
				exclude: node_modules
			}
		],
	},
	vue: {
		loaders: {
			stylus: ExtractTextPlugin.extract('style', 'css!stylus?paths=node_modules/jeet/stylus/'),
			scss: 'style!css!sass',
			exclude: node_modules,
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
