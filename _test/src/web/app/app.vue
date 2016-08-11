<template>
	<div class="wrapper">
		<router-view></router-view>
	</div>
</template>

<style lang="stylus">

	body {
		background: #efefef;
		padding: 20px;
	}

	.wrapper {
		width: 80%;
		min-height: 200px;
		background: #fff;
		margin: 40px auto;
		padding: 40px;
		border-radius: 2px;
		box-shadow: 0px 2px 4px rgba(0,0,0, 0.1);
	}

</style>
