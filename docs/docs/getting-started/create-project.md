---
sidebar_position: 2
---

# Create Your First Project

Now, let's start creating your project.

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

Let's also run:
```
mify add frontend <frontend-service-name>
```
As the name of this command suggests, this will be your project's frontend.
After running this command you will see a `js-services` directory.

## Workspace structure

Let's get a quick overview of the generated structure. The `go-services`
directory which was created by the `add service` command is a place where all
services written in Go being put. If you run the same command with the flag
`--language python` it will generate code in `py-services`. And when you add
frontend it is being put into `js-services`.

There is also a `schemas` directory, which is very important. In order to
support writing services in multiple languages Mify need a way to describe
language-agnostic API description. We use
[OpenAPI](https://spec.openapis.org/oas/latest.html) schemas for that.

The main API file for your service is located at
`schemas/<service-name>/api/api.yaml`. After modifying API you need to call
`mify generate <service-name>` to update the code.

After creating your schema you can start writing your code - the entrypoint for
the generated handlers will be at
`go-services/internal/<service-name>/handlers/<path/to/api>/service.go`. We
suggest putting service logic in `go-services/internal/<service-name>`.

Let's jump to writing code in our guides: [Creating an example project](/docs/guides/overview)

