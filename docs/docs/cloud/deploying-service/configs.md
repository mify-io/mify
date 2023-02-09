---
sidebar_position: 1
---

# Configs

## Config file location

Mify Cloud provides a way to define various additional configuration for service deployment.
It is defined in `schemas/<service-name>/cloud.mify.yaml`.
Everything described here should be put in this file.

## Publishing service

By default your backend service is not available on the web, only inside the cloud.
This is useful when you have some internal services which shouldn't directly interact with users.
But if you want to make backend service public, for instance if you're making public API or frontend outside the Mify Workspace, you need to add `publish: true` to your cloud.mify.yaml.
:::tip
You don't have to manually add `publish` flag if you have frontend inside Mify Workspace.
:::

### Custom domain

If you want to setup custom domain for you service we have `domain` property for that, here's an example for a backend and frontend configuration for that:

`schemas/<backend-service-name>/cloud.mify.yaml`
```yaml
publish: true # should be true if you want custom domain
domain:
    custom_hostname: example.com
    path: /api/ # backend will be available at example.com/api
```

`schemas/<frontend-service-name>/cloud.mify.yaml`
```yaml
publish: true
domain:
    custom_hostname: example.com
```

### Resource limits

By default service will have `200m` for CPU limits and requests and `200M` for
memory. To change that, add resources field:

`schemas/<service-name>/cloud.mify.yaml`
```yaml
resources:
    prod:
        cpu_limit: 400m
        memory_limit: 400M
```
The format of limit quantities is taken from Kubernetes.

### Passing environment variables

Sometimes you will need to pass additional environment variables for your services, for instance for authentication keys,or external API hostnames. Here's what we have for that:

`schemas/<service-name>/cloud.mify.yaml`
```yaml
env_vars:
  # This env variable is just some constant
  SOME_STATIC_BUT_NOT_SECRET_VALUE:
    value: value-of-this-env
  # This could be some secret key, so you wouldn't put raw value here
  SECRET_ENV:
    secret_name: name-of-secret # this is the name of the secret in mify cloud (kubernetes secret)
```

To add secret you need to use this kubectl command:
```
kubectl create secret generic <name-of-secret> --from-file=value="<path/to/api>"
```

:::caution
Remove trailing newline from file with secret content.
:::
