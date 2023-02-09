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
    <span> | </span>
    <a href="https://join.slack.com/t/mifyio/shared_invite/zt-1llnbiio6-lG45E696QOEVzHb0__Qqxw">
      Slack Channel
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
- Multiple language code generation (Right now it's Go and Python for backend, NuxtJS, React on Typescript based frontends)

## Installation

You can get the latest Mify CLI from [Releases](https://github.com/mify-io/mify/releases).

Before using it you need to install and start Docker which is used for running
some code generation tasks, You can refer to Docker's
[guide](https://docs.docker.com/get-docker/) for installation.

### Supported platforms

Right now Mify should work on Linux, Mac and WSL.

### Development prerequisites

At this moment Mify supports Go and Python language based templates for
backends, and NuxtJS and React on Typescript for frontends, here's what you
need to install before starting developing in your choosen template:

- Go:
  - Go >= 1.18

- Python (Beta):
  - Python >= 3.8
  - python3-pip
  - python3-venv

- NuxtJS, React:
  - Node >= 18.12.1
  - Yarn

### Getting the last version

You can always install mify from main branch using Go:
```sh
$ go install github.com/mify-io/mify/cmd/mify@latest
```

## Quick Start

Create your first project: https://mify.io/docs

Guide on how to create sample backend and frontend app: https://mify.io/docs/guides/overview

If you have any questions, join our [Slack channel](https://join.slack.com/t/mifyio/shared_invite/zt-1llnbiio6-lG45E696QOEVzHb0__Qqxw)!

## License
[Apache 2.0](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0))
