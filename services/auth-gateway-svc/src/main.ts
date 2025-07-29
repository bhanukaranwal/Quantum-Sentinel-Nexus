// services/auth-gateway-svc/src/main.ts
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { createProxyMiddleware } from 'http-proxy-middleware';
import { INestApplication } from '@nestjs/common';

// --- Configuration ---
const INGESTION_SVC_URL =
  process.env.INGESTION_SVC_URL || 'http://ingestion-svc:8080';
const PORT = process.env.PORT || 3001;

/**
 * Sets up a proxy middleware for a specific route.
 * @param app The NestJS application instance.
 * @param context The route context (e.g., '/ingest').
 * @param target The downstream service URL to proxy to.
 */
function setupProxy(app: INestApplication, context: string, target: string) {
  const proxy = createProxyMiddleware({
    target,
    changeOrigin: true,
    pathRewrite: {
      [`^/api${context}`]: '', // Rewrite /api/ingest to /
    },
    onProxyReq: (proxyReq, req, res) => {
      console.log(
        `[Proxy] Forwarding request from ${req.url} to ${target}${proxyReq.path}`,
      );
    },
    onError: (err, req, res) => {
      console.error('[Proxy] Error:', err);
      res.writeHead(500, {
        'Content-Type': 'text/plain',
      });
      res.end('Proxy error: Could not connect to downstream service.');
    },
  });
  app.use(`/api${context}`, proxy);
}

/**
 * The bootstrap function to initialize and start the NestJS application.
 */
async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  // Enable Cross-Origin Resource Sharing (CORS) for the frontend.
  app.enableCors();

  // --- Middleware & Proxies ---
  // In a real-world scenario, this proxying logic would be more sophisticated,
  // likely managed by a dedicated API gateway like Envoy, Kong, or Traefik,
  // configured via Kubernetes CRDs. This implementation provides a basic
  // application-level proxy for development purposes.

  // TODO: Implement OIDC authentication middleware (e.g., using Passport.js and passport-jwt).
  // TODO: Implement mTLS validation middleware.
  // TODO: Implement rate-limiting middleware (e.g., using @nestjs/throttler).

  // Proxy requests for '/api/ingest' to the 'ingestion-svc'.
  setupProxy(app, '/ingest', INGESTION_SVC_URL);

  // You would add more proxies here for other services.
  // setupProxy(app, '/scoring', SCORING_SVC_URL);

  // --- Server Start ---
  await app.listen(PORT);
  console.log(`🚀 Auth Gateway Service is running on: http://localhost:${PORT}`);
  console.log(`  -> Proxying /api/ingest to ${INGESTION_SVC_URL}`);
}

bootstrap();
