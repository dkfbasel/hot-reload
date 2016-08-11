<style lang="stylus" scoped>

	// test jeet integration
	@require('jeet/_jeet.styl');

	button {
		padding: 8px 12px;
		background: #efefef;
		border: 1px solid tint(#000, 60%);

		&:hover {
			background: darken(#efefef, 10%);
		}
	}

	pre {
		background #efefef;
		margin-top: 20px;
		min-height: 80px;
		line-height: 1.3em;
		display: block;
		color: #000;
		padding: 20px;
	}

</style>

<template>
	<div>
		<h1>This is a first test for webpack live reload in a docker container</h1>
		<p>In this test, we will try to {{ content }}</p>

		<button v-on:click="requestApiData">Request data from api</button>

		<pre>{{ console }}</pre>
	</div>
</template>

<script lang="babel">

	module.exports = {
		data: () => {
			return {
				content: 'use webpack with live reload in a docker container',
				console: ''
			};
		},
		methods: {
			requestApiData: function() {

				this.log('--');
				this.log('start request from api');

				// make api request
				this.$http.get('/api').then((response) => {

					this.log('api response:');
					this.log('- ' + response.data);

				}, (response) => {
					this.log('there was an error');
					this.log(response);
				});

			},

			log: function(message) {
				this.console += message + '\n';
			}
		}
	};

</script>
