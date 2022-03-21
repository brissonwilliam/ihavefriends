const { createProxyMiddleware } = require('http-proxy-middleware');

const BACKEND_HOST = 'http' + process.env.REACT_APP_TLS + '://' + process.env.REACT_APP_BACKEND_HOST;

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: BACKEND_HOST,
      changeOrigin: true,
    })
  );
};