<script>
export default {
	data() {
		return {
			globalData: {
				toast: Object,//全局提示，必须在页面初始化
			},
		}
	},
	onLaunch: function () {
		console.log('App Launch')
	},
	onShow: function () {
		console.log('App Show')
	},
	onHide: function () {
		console.log('App Hide')
	}
}
</script>

<style lang="scss">
@import url("style/king_style.css");
/*每个页面公共css */
body{
	background: #f7f7f7;
}
</style>
