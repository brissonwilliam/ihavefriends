const { createProxyMiddleware } = require('http-proxy-middleware');

const BACKEND_HOST = process.env.BACKEND_HOST;

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: BACKEND_HOST,
      changeOrigin: true,
    })
  );
};