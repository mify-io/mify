FROM node:lts-alpine as builder

WORKDIR /app

{{- range .Clients}}
ARG {{.EndpointEnvName}}
ENV {{.EndpointEnvName}}=${
{{- .EndpointEnvName -}}
}
{{- end}}

COPY package.json tsconfig.json ./
COPY public/ public/
COPY src/ src/

RUN yarn && yarn build

FROM nginx:1.25.1-alpine

COPY <<EOF /etc/nginx/conf.d/default.conf
RUN cat server_tokens off;
server {
    listen       80;
    server_name  localhost;
    location / {
        root   /usr/share/nginx/html;
        index  index.html index.htm;
        try_files $uri /index.html;
    }
}
EOF

COPY --from=builder /app/build /usr/share/nginx/html

RUN touch /var/run/nginx.pid &&  \
    chown -R nginx:nginx /var/run/nginx.pid /usr/share/nginx/html /var/cache/nginx /var/log/nginx /etc/nginx/conf.d
USER nginx
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
