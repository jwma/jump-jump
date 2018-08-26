import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
const flyio = require('flyio')

Vue.config.productionTip = false

flyio.interceptors.request.use((request) => {
  request.headers['Content-Type'] = 'application/x-www-form-urlencoded'
  return request
})
Vue.prototype.$http = flyio

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
