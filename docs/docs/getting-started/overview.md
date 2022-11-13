---
sidebar_position: 1
---

# Installation

You can get the latest Mify CLI from [GitHub](https://github.com/mify-io/mify/releases).

Before using it you need to install and start Docker which is used for running
some code generation tasks, You can refer to Docker's
[guide](https://docs.docker.com/get-docker/) for installation.

## Supported platforms

Right now Mify should work on Linux, Mac and WSL.

## Development prerequisites

At this moment Mify supports Go and Python language based templates for backends,
and NuxtJS for frontends, here's what you need to install before starting developing
in your choosen template:

- Go:
  - Go >= 1.18

- Python:
  - Python >= 3.8
  - python3-pip
  - python3-venv

- NuxtJS:
  - Node >= 18.12.1
  - Yarn

## Getting the last version

You can always install mify from main branch using Go:
```sh
$ go install github.com/mify-io/mify/cmd/mify@latest
```
