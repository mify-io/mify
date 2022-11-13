{
  "name": "{{.ServiceName}}",
  "version": "1.0.0",
  "scripts": {
    "dev": "NODE_OPTIONS=--openssl-legacy-provider nuxt",
    "build": "NODE_OPTIONS=--openssl-legacy-provider nuxt build",
    "generate": "NODE_OPTIONS=--openssl-legacy-provider nuxt generate",
    "start": "NODE_OPTIONS=--openssl-legacy-provider nuxt start"
  },
  "dependencies": {
    "nuxt": "^2.15.8",
    "superagent": "^7.1.1"
  }
}
