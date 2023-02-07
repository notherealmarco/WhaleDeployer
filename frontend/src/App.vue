<script>
export default {
	props: ["user_id", "name", "date", "comments", "likes", "photo_id", "liked"],
	data: function () {
		return {
			theme: "dark",
		}
	},
	methods: {
		setDarkTheme(dark) {
			this.theme = dark ? "dark" : "light"
		},
		switchTheme() {
			this.theme = this.theme === "light" ? "dark" : "light"
		},
	},

	// Called when the root view is mounted
	mounted() {
		//this.theme = this.$grammy.colorScheme
	}
}
</script>

<template>
	<v-app :theme="theme">
		<v-container class="fill-height">
			<v-row align-content="center" justify="center" class="fill-height">
				<router-view></router-view>
				<div v-if="theme === 'light'" style="color: gray">why would you :/</div>
			</v-row>
		</v-container>
	</v-app>
</template>

<style>
</style>