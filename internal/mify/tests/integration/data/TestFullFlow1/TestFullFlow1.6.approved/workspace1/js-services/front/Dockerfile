FROM node:lts as builder

WORKDIR /app

COPY . .

WORKDIR /app/front

RUN yarn install \
  --prefer-offline \
  --frozen-lockfile \
  --non-interactive \
  --production=false

RUN yarn build

RUN rm -rf node_modules && \
  NODE_ENV=production yarn install \
  --prefer-offline \
  --pure-lockfile \
  --non-interactive \
  --production=true

FROM node:lts

WORKDIR /app

COPY --from=builder /app  .

WORKDIR /app/front

ENV HOST 0.0.0.0
ENV PORT 80
EXPOSE 80

CMD [ "yarn", "start" ]
