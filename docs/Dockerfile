FROM node:lts-alpine as builder

WORKDIR /app

COPY ./mify-docs .

WORKDIR /app/mify-docs

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

FROM node:lts-alpine

WORKDIR /app

COPY --from=builder /app  .

WORKDIR /app/mify-docs

ENV HOST 0.0.0.0
ENV PORT 80
EXPOSE 80

CMD [ "yarn", "start" ]
