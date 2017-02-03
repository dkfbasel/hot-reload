var path = require('path');
var webpack = require('webpack');
var poststylus = require('poststylus');
var autoprefixer = require('autoprefixer');

// note: we prefer using includes over excludes, as this will give us finer
// control over what is actually transpiled
var appDirectory = path.resolve(__dirname, 'app');
var includes = [appDirectory];

module.exports = {
	devServer: {
		historyApiFallback: true,
		noInfo: true,
		contentBase: path.resolve(__dirname, 'build'),
		host: '0.0.0.0',
		port: 3000
	},
	performance: {
		hints: false
	},
	devtool: '#eval-source-map',
	entry: {
		app: [path.resolve(__dirname, 'app/main.js')]
	},
	output: {
		path: path.resolve(__dirname, 'build/assets'),
		filename: '[name].bundle.js',
		publicPath: '/assets/'
	},
	module: {
		rules: [
			{
				// parse vue components
				test: /\.vue$/,
				loader: 'vue-loader',
				include: includes
			}, {
				// parse css styles
				test: /\.css$/,
				use: ['style-loader','css-loader','postcss-loader'],
				include: includes
			}, {
				// parse javascript files
				test: /\.js$/,
				loader: 'babel-loader',
				query: {
					presets: ['es2015', 'stage-0'],
					plugins: ['transform-runtime']
				},
				include: includes
			}, {
				// parse stylus styles
				test: /\.styl$/,
				use: [
					{loader: 'style-loader'},
					{loader: 'css-loader'},
					{
						loader:'stylus-loader',
						options: {
							ident: 'stylus',
							use: [
								poststylus([
									autoprefixer({
										browsers: ['iOS >= 6', 'ie >= 9']
									})
								])
							]
						}
					}
				],
				include: includes
			}
		]
	},
	resolve: {
		alias: {
			vue: 'vue/dist/vue.js'
		}
	}
};


if (process.env.NODE_ENV === 'production') {
	module.exports.devtool = '#source-map';

	// http://vue-loader.vuejs.org/en/workflow/production.html
	module.exports.plugins = (module.exports.plugins || []).concat([
		new webpack.DefinePlugin({
			'process.env': {
				NODE_ENV: '"production"'
			}
		}),
		new webpack.optimize.UglifyJsPlugin({
			sourceMap: true,
			compress: {
				warnings: false
			}
		}),
		new webpack.LoaderOptionsPlugin({
			minimize: true
		})
	]);
}
