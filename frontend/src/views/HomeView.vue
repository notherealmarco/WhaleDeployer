<script>
export default {
	data: function () {
		return {
			some_data: [],
			show_plus: true,

			show_dialog: false,
			dialog_code: "",
			dialog_text: "",
		}
	},
	methods: {
		async refresh() {
			try {
				let response = await this.$axios.get("/projects")

				response.data.forEach(i => {
					if (i.deploy_key === true) i.deploy_key = '1'
					else i.deploy_key = '0'

					i.is_new = false
					i.expanded = false
				});

				this.some_data = response.data
				this.show_plus = true
			} catch (e) {
				alert(e.toString())
			}
		},

		async ldelete(name) {
			try {
				let response = await this.$axios.delete("/projects/" + name)
				this.refresh()
			} catch (e) {
				alert(e.toString())
			}
		},

		async build(name) {
			this.$axios.post("/projects/" + name).then((r) => {
				this.refresh()
			}).catch((e) => {
				this.refresh()
			})
			this.refresh()
		},
		async logs(name) {
			this.$axios.get("/projects/" + name + "/logs").then((r) => {
				this.dialog_text = name + " logs:"
				this.dialog_code = r.data
			}).catch((e) => {
				this.refresh()
			})
		},
	
		async show_refresh_logs(name) {
			this.show_dialog = true
			while(this.show_dialog) {
				await this.logs(name)
				await this.sleep(1)
			}
		},
		async sleep(seconds) {
			return new Promise((resolve) => setTimeout(resolve, seconds * 1000));
		},

		async updateInfo(item) {
			await this.$axios.post("/projects", {
				name: item.name,
				description: item.description,
				path: item.path,
				git_url: item.git_url,
				git_branch: item.git_branch,
				dockerfile: item.dockerfile,
				image_name: item.image_name,
				image_tag: item.image_tag,
				deploy_key: item.deploy_key === '1' ? true : false,
			}).then((r) => {
				this.refresh()
				if (item.deploy_key) {
					this.dialog_text = "Your Git deploy key:"
					this.dialog_code = r.data
					this.show_dialog = true
				}

			}).catch((e) => {
				alert(e.toString())
				console.log(e)
				this.refresh()
			})
		},
		async deleteProject(name) {
			await this.$axios.delete("/projects/" + name).then((r) => {
				this.refresh()
			}).catch((e) => {
				alert(e.toString())
				console.log(e)
				this.refresh()
			})
		},
		async newItem() {
			let item = {
				name: "new project",
				description: "new project",
				path: "",
				git_url: "",
				git_branch: "",
				dockerfile: "",
				image_name: "",
				image_tag: "",
				deploy_key: false,
				is_new: true,
				expanded: true,
			}
			this.some_data.unshift(item)
			this.show_plus = false
		}
	},
	mounted() {
		this.refresh();

		//console.log(this.$vuetify)

		//this.$grammy.MainButton.setParams({
		//	text: "pulsantone inutile",
		//	color: this.$grammy.themeParams.hint_color,
		//	is_active: false,
		//	is_visible: false,
		//});
	}
}
</script>

<template>

	<div class="text-center">
		<v-dialog v-model="show_dialog">
			<v-card>
				<v-card-text style="text-align: center">
					<b>{{ dialog_text }}</b><br /><pre>{{ dialog_code }}</pre>
				</v-card-text>
				<v-card-actions>
					<v-btn color="primary" block @click="show_dialog = false">Close Dialog</v-btn>
				</v-card-actions>
			</v-card>
		</v-dialog>
	</div>

	<v-container>
		<v-row>
			<v-col cols="12" sm="1" md="2" lg="3"></v-col>
			<v-col cols="12" sm="10" md="8" lg="6">

				<div class="d-flex flex-row-reverse">
					<v-btn v-if="show_plus" class="d-flex mb-3 ml-2" icon="mdi-plus" color="success" size="2.7em" @click="newItem()"></v-btn>

					<v-btn class="d-flex mb-3" icon="mdi-weather-night" color="amber" size="2.7em" @click="$root.switchTheme()"></v-btn>
				</div>

				<v-card v-for="i in some_data" class="mx-auto mb-5">
					<!--<v-img src="https://cdn.vuetifyjs.com/images/cards/sunshine.jpg" height="200px" cover></v-img>-->

					<v-row>
						<v-col cols="10">
							<v-card-title>
								{{ i.name }}
							</v-card-title>

							<v-card-subtitle>
								{{ i.description }}
							</v-card-subtitle>
						</v-col>

						<v-col cols="2" class="d-flex align-end flex-column">
							<div class="d-flex mt-3 me-4">
								<v-progress-circular indeterminate color="amber"
									v-if="i.status === 'running'"></v-progress-circular>
								<v-icon v-if="i.status === 'success'" color="green" size="x-large"
									style="font-size: 2.5em">mdi-check-circle</v-icon>
								<v-icon v-if="i.status === 'fail'" color="red" size="x-large"
									style="font-size: 2.5em">mdi-alert-circle</v-icon>
								<v-icon v-if="i.status === 'unknown'" color="grey" size="x-large"
									style="font-size: 2.5em">mdi-help-circle</v-icon>
							</div>
						</v-col>
					</v-row>

					<v-card-actions>
						<v-btn color="primary" variant="text" @click="build(i.name)">
							<span v-if="i.state !== 'unknown'">Re-Deploy</span>
							<span v-if="i.state === 'unknown'">First Deploy</span>
						</v-btn>

						<v-btn color="secondary" variant="text" @click="show_refresh_logs(i.name)">
							Logs
						</v-btn>

						<v-spacer></v-spacer>

						<v-btn :icon="i.expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'"
							@click="i.expanded = !i.expanded"></v-btn>
					</v-card-actions>

					<v-expand-transition>
						<div v-show="i.expanded">
							<v-divider></v-divider>

							<v-card-text>
								<form @submit.prevent="submit">
									<v-text-field v-model="i.name" label="Name" v-if="i.is_new"></v-text-field>

									<v-text-field v-model="i.description" label="Description"></v-text-field>

									<v-text-field v-model="i.git_url"
										label="Git URI (only use SSH for private repos)"></v-text-field>

									<v-row>
										<v-col>
											<v-text-field v-model="i.git_branch" label="Git branch"></v-text-field>
										</v-col>
										<v-col>
											<v-text-field v-model="i.path" label="File system path"></v-text-field>
										</v-col>
									</v-row>

									<v-text-field v-model="i.dockerfile" label="Dockerfile name"></v-text-field>

									<v-row>
										<v-col cols="6"><v-text-field v-model="i.image_name"
												label="Image name"></v-text-field></v-col>

										<v-col cols="6"><v-text-field v-model="i.image_tag"
												label="Image tag"></v-text-field></v-col>
									</v-row>

									<v-checkbox value="1" label="I have a private repository" type="checkbox"
										color="primary" v-model="i.deploy_key"></v-checkbox>

									<v-btn class="me-4" color="primary" type="submit" @click="updateInfo(i)">
										Save
									</v-btn>

									<v-btn @click="deleteProject(i.name)" color="error" variant="text">
										Delete
									</v-btn>
								</form>
							</v-card-text>
						</div>
					</v-expand-transition>
				</v-card>


			</v-col>
			<v-col xs="12" sm="1" md="2" lg="3"></v-col>
		</v-row>
	</v-container>


</template>

<style>

</style>
