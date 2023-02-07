import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vuetify, { transformAssetUrls } from 'vite-plugin-vuetify'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode, ssrBuild }) => {
	const ret = {
		plugins: [
			vue({
				template: { transformAssetUrls }
			}),
			vuetify({ autoImport: true }),
		],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url))
			}
		},
	};
	if (command === 'serve') {
		ret.define = {
			"__API_URL__": JSON.stringify("http://localhost:3000/api"),
		};
	} else {
		ret.define = {
			"__API_URL__": JSON.stringify("/api"),
		};
	}
	return ret;
})
