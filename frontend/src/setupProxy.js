const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'https://localhost:8443',
      changeOrigin: true,
      secure: false, // Игнорируем ошибки самоподписанного сертификата
      pathRewrite: {
        '^/api': '/api', // Оставляем путь как есть
      },
      onProxyReq: (proxyReq, req, res) => {
        // Можно добавить логирование или заголовки
        console.log(`Proxy: ${req.method} ${req.path} -> ${proxyReq.path}`);
      },
      onError: (err, req, res) => {
        console.error('Proxy error:', err);
        res.status(500).json({ error: 'Proxy error' });
      }
    })
  );
};