<template>
  <div class="home">
    <div class="form">
      <h3 class="form-title">创建短链接</h3>
      <div class="tips" v-show="tips != ''">{{ tips }}</div>
      <div class="form-item">
        <div class="form-label"><label for="url">URL</label></div>
        <div class="form-field"><input id="url" type="url" placeholder="-" v-model="form.url"></div>
      </div>
      <div class="form-item">
        <div class="form-label"><label for="description">描述</label></div>
        <div class="form-field"><input id="description" type="text" placeholder="-" v-model="form.description"></div>
      </div>
      <div class="form-item button-wrapper">
        <button class="btn" @click="submit">提 交</button>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'

export default {
  name: 'home',
  data () {
    return {
      tips: null,
      form: {
        url: null,
        description: null
      }
    }
  },
  methods: {
    submit () {
      const url = this.form.url, description = this.form.description, isEnabled = true
      if (!url) {
        this.tips = '请填写链接'
        return
      }
      Vue.prototype.$http.post(`${process.env.VUE_APP_API_ADDR}/admin/link`, {url, description, isEnabled})
        .then(response => {
          if (response.data.code === 4999) {
            this.tips = response.data.msg
          } else {
            this.tips = `短链接创建成功，ID：${response.data.slug}`
            this.form.url = null
            this.form.description = null
          }
        })
    }
  }
}
</script>

<style scoped>
.home {
  font-size: 18px;
  text-align: center;
}
</style>
