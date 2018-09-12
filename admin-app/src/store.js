import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    token: null,
    username: null
  },
  mutations: {
    login(state, payload) {
      state.token = payload.token
      state.username = payload.username
    }
  },
  actions: {
    login({ commit, state }, formData) {
      return Vue.prototype.$http.post(`${process.env.VUE_APP_API_ADDR}/login/`, formData)
        .then(response => {
          const data = response.data
          if (data.code === 0) {
            // 缓存数据
            window.localStorage.setItem('token', data.token)
            window.localStorage.setItem('username', data.username)
            
            // 更新state
            commit({ type: 'login', token: data.token, username: data.username })
          }

          return response.data
        })
    }
  }
})
