---
sidebar_position: 1
slug: /
---

# Get started

Mify is an open-source developer tool which generates and maintans
infrastructure code for your cloud service. It provides infra code for stuff
like logs, metrics, and API, allowing engineers to focus on meaningful code.
Services generated by Mify can be deployed in our cloud without any tedious
configuration.

## Installing Mify

You can get the latest Mify CLI from [GitHub](https://github.com/mify-io/mify/releases).

Before using it you need to install and start Docker which is used for running
some code generation tasks, You can refer to Docker's
[guide](https://docs.docker.com/get-docker/) for installation.
Right now Mify should work on Linux, Mac and WSL.

## Creating your first project

After installing Mify CLI call this command:
```
mify init <project-name>
```

You will see a directory called `<project-name>` which is now your workspace where you can add your services.
Go into workpace with `cd <project-name>` and create your first service:
```
mify add service <service-name>
```
You will see that this command generated a service in `go-services` directory.
In `schemas/<service-name>/api.yaml` there is OpenAPI schema which is a starting point
for creating service handlers. Run `mify generate` after changing the schema and
new handler will appear at `go-services/internal/<service-name>/handlers/path/to/api/service.go`.

And if you also run:
```
mify add frontend <frontend-service-name>
```
As the name of this command suggests, this will be your project's frontend.
After running this command you will find it in a `js-services` directory.

### Running the service

Get into go-services directory and install dependencies:

```
$ cd go-services
$ go mod tidy
```

Build and run the service:

```
$ go run ./cmd/<service-name>
```

Follow this [guide](/docs/cloud/overview)  to deploy it to Mify Cloud with single command `mify cloud deploy`.


:::tip
Check out our [guides](/docs/guides/overview) to get more undestanding on how to
write services with Mify.
:::

##

### Development prerequisites

At this moment Mify supports Go and Python language based templates for
backends, and NuxtJS and React on Typescript for frontends, here's what you
need to install before starting developing in your choosen template:

- Go:
  - Go >= 1.18

- Python (beta):
  - Python >= 3.8
  - python3-pip
  - python3-venv

- NuxtJS, React:
  - Node >= 18.12.1
  - Yarn

## Getting the last version

You can always install mify from main branch using Go:
```sh
$ go install github.com/mify-io/mify/cmd/mify@latest
```
## Community

Join our [Slack
channel](https://join.slack.com/t/mifyio/shared_invite/zt-1llnbiio6-lG45E696QOEVzHb0__Qqxw)
if have suggestions or if you need any help or just to talk about cloud services.
