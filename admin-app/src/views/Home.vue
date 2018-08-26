<template>
  <div class="home">
    <div class="form">
      <h3 class="form-title">创建短链接</h3>
      <div class="tips" v-show="tips != ''">{{ tips }}</div>
      <div class="form-item">
        <div class="form-label"><label for="url">链接</label></div>
        <div class="form-field"><input id="url" type="url" v-model="form.url"></div>
      </div>
      <div class="form-item">
        <div class="form-label"><label for="description">描述</label></div>
        <div class="form-field"><input id="description" type="text" v-model="form.description"></div>
      </div>
      <div class="form-item button-wrapper">
        <button class="submit-btn" @click="submit">提 交</button>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'

export default {
  name: 'home',
  data() {
    return {
      tips: null,
      form: {
        url: null,
        description: null
      }
    }
  },
  methods: {
    submit() {
      const url = this.form.url, description = this.form.description, isEnabled = true
      if (!url) {
        this.tips = '请填写链接'
        return
      }
      Vue.prototype.$http.post('/api/admin/link', {url, description, isEnabled})
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

.form {
  margin-top: 30px;
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

.submit-btn {
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

