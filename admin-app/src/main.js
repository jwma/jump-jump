import Vue from 'vue'
import App from './App.vue'     // 应用的入口单文件组件
import router from './router'   // 应用的路由实例，路由配置所在的地方
import store from './store'     // 应用的状态仓库实例
const flyio = require('flyio')

Vue.config.productionTip = false

/**
 * 使用 flyio 的拦截器设置每次请求的 Content-Type=application/x-www-form-urlencoded
 * 这样后端处理请求体中的参数时就会方便些
 */
flyio.interceptors.request.use((request) => {
  request.headers['Content-Type'] = 'application/x-www-form-urlencoded'
  request.headers['Authorization'] = store.state.token
  return request
})
Vue.prototype.$http = flyio

// 如果访问需要登录态的路由，则检查是否登录
router.beforeEach((to, form, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!store.state.token) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    next()
  }
})

// 首次加载尝试从 localStorage 获取登录信息缓存
store.commit({
  type: 'login',
  token: window.localStorage.getItem('token'),
  username: window.localStorage.getItem('username')
})

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')

