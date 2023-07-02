import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import router from './router'
import VueAxios from 'vue-axios';

Vue.config.productionTip = false

const client = axios.create({
  baseURL: '/api',
})
Vue.use(VueAxios, client);

new Vue({
  router,
  render: function (h) { return h(App) }
}).$mount('#app')
