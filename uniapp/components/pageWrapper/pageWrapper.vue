<template>
	<view :style="{ height: windowHeight + 'px' }">
		<slot></slot>
	</view>
</template>

<script>
export default {
	name: "pageWrapper",
	data() {
		return {
			windowHeight: 0,
		};
	},
	props: {
		subHeight: {
			type: Number,
			default: 0,
		}
	},
	mounted() {
		uni.getSystemInfo({
			success: (res) => {
				this.windowHeight = res.windowHeight - this.subHeight;
			},
		});
	}
}
</script>

<style lang="scss"></style>
