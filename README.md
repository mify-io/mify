<h1 align="center">
  <br>
  <a href="https://mify.io"><img src="https://raw.githubusercontent.com/mify-io/mify/main/docs/static/img/logo.png" alt="Mify" width="200"></a>
</h1>

<div align="center">
</div>
<div align="center">
  <strong>Microservice infrastructure for you</strong>
</div>
<div align="center">
  A code generation tool to help you build cloud backend services
</div>

<div align="center">
</div>

<div align="center">
  <h3>
    <a href="https://mify.io">
      Website
    </a>
    <span> | </span>
    <a href="https://mify.io/docs">
      Docs
    </a>
    <span> | </span>
    <a href="https://github.com/mify-io/mify/blob/main/.github/CONTRIBUTING.md">
      Contributing
    </a>
  </h3>
</div>

<div align="center">
</div>

[![Go](https://github.com/mify-io/mify/actions/workflows/go.yml/badge.svg)](https://github.com/mify-io/mify/actions/workflows/go.yml)

## Features
- OpenAPI http server generation
- Built-in Prometheus metrics
- Structured logging
- Multiple language code generation (Right now it's Go for backend, and NuxtJS based frontend)

## Installation

Install Docker, it is needed for OpenAPI generation.

Then install mify CLI using Go:
```sh
$ go install github.com/mify-io/mify/cmd/mify@latest
```
## Quick Start

Create your first project: https://mify.io/docs/getting-started/create-project

Guide on how to create sample backend and frontend app: https://mify.io/docs/guides/create-service

## License
[Apache 2.0](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0))
