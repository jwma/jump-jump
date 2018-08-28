<template>
  <div id="app">
    <div id="nav">
      <router-link to="/">Jump-Jump</router-link> |
      <router-link to="/list">Editor</router-link>
    </div>
    <router-view/>
  </div>
</template>

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
