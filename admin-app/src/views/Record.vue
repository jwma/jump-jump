<template>
  <div class="record">
    <div class="search-wrapper">
      <div><input type="text" v-model="slug" placeholder="请输入短链接ID"> <button class="btn" @click="search">查 询</button></div>
      <div style="margin-top: 10px;"></div>
    </div>
    <div class="form" v-show="showForm">
      <div class="tips" v-show="tips != ''">{{ tips }}</div>
      <div class="form-item">
        <div class="form-label"><label for="url">URL</label></div>
        <div class="form-field"><input type="url" v-model="form.url"></div>
      </div>
      <div class="form-item">
        <div class="form-label"><label for="description">描述</label></div>
        <div class="form-field"><input type="url" v-model="form.description"></div>
      </div>
      <div class="form-item">
        <div class="form-label"><label><input type="checkbox" v-model="form.isEnabled"> 是否开启</label></div>
      </div>
      <div class="form-item button-wrapper">
        <button class="btn" @click="submit">更 新</button>
      </div>
    </div>
  </div>
</template>
<script>
import Vue from 'vue'

export default {
  name: 'Record',
  data () {
    return {
      slug: null,
      tips: null,
      showForm: false,
      form: {
        url: null,
        description: null,
        isEnabled: false
      }
    }
  },
  methods: {
    search () {
      Vue.prototype.$http.get(`/api/admin/link?slug=${this.slug}`)
        .then(response => {
          if (response.data.code === 4999) {
            alert(response.data.msg)
          } else {
            const link = response.data.link
            this.form.url = link.Url,
            this.form.description = link.Description
            this.form.isEnabled = link.IsEnabled
            this.showForm = true
          }
        })
    },
    submit() {
      const url = this.form.url, description = this.form.description, isEnabled = true
      if (!url) {
        this.tips = '请填写链接'
        return
      }

      Vue.prototype.$http.patch(`/api/admin/link?slug=${this.slug}`, {url, description, isEnabled})
        .then(response => {
          if (response.data.code === 4999) {
            this.tips = response.data.msg
          } else {
            this.tips = '更新成功'
          }
        })
    }
  }
}
</script>

<style scoped>
.record {
  text-align: center;
}

.search-wrapper {
  margin: 0 auto;
  width: 300px;
}

.search-wrapper input {
  padding: 5px 3px;
  text-align: center;
  width: 200px;
}
</style>

