# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY web/package*.json ./

RUN npm ci --omit=dev

COPY web ./

RUN npm run build

# Production stage  
FROM node:20-alpine

WORKDIR /app

ENV NODE_ENV=production

# Copy only production dependencies and build output
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/src ./src

EXPOSE 3000

CMD ["npm", "start"]