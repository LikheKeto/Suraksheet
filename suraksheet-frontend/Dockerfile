FROM node:20-slim AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

ENV NODE_OPTIONS="--max-old-space-size=8000"

RUN npm run build

FROM node:20-slim AS runner

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install --omit=dev

COPY --from=builder /app/.svelte-kit /app/.svelte-kit
COPY --from=builder /app/static /app/static
COPY --from=builder /app/build /app/build
COPY --from=builder /app/src /app/src

ENV NODE_ENV=production

EXPOSE 3000

CMD ["node", "build"]