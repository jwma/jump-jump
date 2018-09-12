<template>
  <!-- 整个应用的根节点 -->
  <div id="app">  
    <div id="nav" v-show="$store.state.token">
      <!-- 导航栏部分定义了两个路由链接 -->
      <router-link to="/">Jump-Jump</router-link> |
      <router-link to="/list">Editor</router-link> |
      {{$store.state.username}} |
      <a href="javascript:" @click="logout">登出</a>
    </div>
    <!-- 路由对应的组件会被渲染到 router-view 这里 -->
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'app',
  methods: {
    logout() {
      this.$store.dispatch('logout')
      this.$router.replace({name: 'login'})
    }
  }
}
</script>


<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
#nav {
  text-align: center;
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}

.form {
  padding: 15px;
  width: 300px;
  margin: 30px auto 0;
  background-color: rgba(238, 238, 238, 0.3);
  border-radius: 5px;
}

.form-title {
  color: #555;
}

.form-item + .form-item {
  margin-top: 15px;
}

.form-field > input {
  padding: 5px 3px;
  text-align: center;
  width: 200px;
}

.button-wrapper {
  margin-top: 15px;
}

.btn {
  font-size: 16px;
  border: 1px solid #999;
  background: #fff;
  outline: none;
  cursor: pointer;
}

.tips {
  font-size: 16px;
  color: crimson;
  margin-bottom: 5px;
}
</style>
