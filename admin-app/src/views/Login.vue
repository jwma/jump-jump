<template>
  <div class="login">
    
    <div class="form">
      <!-- 所有提示信息在这里显示 -->
      <div class="tips" v-show="tips != ''">{{ tips }}</div>

      <div class="form-item">
        <div class="form-label"><label for="username">用户名</label></div>
        <div class="form-field"><input id="username" type="text" placeholder="" v-model="form.username"></div>
      </div>
      <div class="form-item">
        <div class="form-label"><label for="password">密 码</label></div>
        <div class="form-field"><input id="password" type="password" placeholder="" v-model="form.password"></div>
      </div>

      <div class="form-item button-wrapper">
        <!-- 绑定按钮点击事件 -->
        <button class="btn" @click="submit">登 录</button>
      </div>
    </div>

  </div>
</template>
<script>
import Vue from 'vue'
import { mapMutations } from 'vuex'

export default {
  name: 'login',
  data () {
    return {
      tips: null,
      form: {
        username: null,
        password: null
      }
    }
  },
  methods: {
    submit() {
      if (!this.form.username || !this.form.password) {
        this.tips = '请输入登录信息'
        return
      }
      this.$store.dispatch('login', this.form).then(rs => {
        // 处理登录成功或失败的结果
        if (rs.code === 0) {
          const redirect = this.$route.query.redirect == '/login' ? '/' : this.$route.query.redirect
          this.$router.push({ path: redirect || '/' })
        } else if (rs.code === 4999) {
          this.tips = rs.msg
        }
      })
    }
  }
}
</script>
<style scoped>
.login {
  font-size: 18px;
  text-align: center;
}
</style>