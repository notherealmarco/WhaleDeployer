import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';

import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'

//import { WebApp } from "@grammyjs/web-app";

import './assets/dashboard.css'
import './assets/main.css'

import vuetify from './plugins/vuetify'


const app = createApp(App)

app.config.globalProperties.$axios = axios;
//app.config.globalProperties.$grammy = WebApp;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);

app.use(vuetify)

app.use(router)
app.mount('#app')
