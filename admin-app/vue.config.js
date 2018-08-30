module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: 'http://web:8081',
        changeOrigin: true,
        pathRewrite: { '^/api': '' }
      }
    }
  }
}