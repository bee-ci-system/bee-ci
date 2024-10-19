# Inspired by https://github.com/vercel/next.js/blob/canary/examples/with-docker-multi-env/docker/production/Dockerfile

FROM node:20-alpine3.20 AS builder

ENV NEXT_TELEMETRY_DISABLED=1

RUN corepack enable pnpm

WORKDIR /app

COPY package.json pnpm-lock.yaml ./

RUN pnpm install

COPY . .

# As of now, this Dockerfile is only used for running locally with Docker. So localhost is fine.

# This is OK. GitHub App Client ID is public information, not an actual secret.
RUN echo "NEXT_PUBLIC_GITHUB_APP_CLIENT_ID=Iv23liiZSvMGEpgOlexa" >> .env.local
# RUN echo "NEXT_PUBLIC_API_BASE_URL=https://bee-ci.karolak.cc/backend/api" >> .env.local
RUN echo "NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api" >> .env.local

RUN pnpm build

FROM node:20-alpine3.20 AS runtime

WORKDIR /app

ENV NODE_ENV=production

RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000

CMD [ "node", "server.js" ]
