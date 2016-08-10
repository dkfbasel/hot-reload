<style lang="stylus" scoped>

	// test jeet integration
	// @require('jeet');

</style>

<template>
	<div>
		<h1>This is a first test</h1>
		<p>In this test, we will try to {{ content }}</p>
	</div>
</template>

<script lang="babel">

	module.exports = {
		data: () => {
			return {
				content: 'use webpack with live reload in a docker container'
			};
		}
	};

</script>
