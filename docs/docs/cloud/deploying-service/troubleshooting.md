---
sidebar_position: 2
---

# Troubleshooting

To check problems with deployed services at the moment you need to use kubectl.
After running `mify cloud init` your `~/.kube/config` will be updated with credentials to
connect to kubernetes. On this page we will list some commands, which will help during troubleshooting.

## Kubernetes Cheatsheet

List pods: `kubectl get pods`

List deployments: `kubectl get deployments`

Get deployment description (list of secrets): `kubectl describe deployment <service-name>`

Follow service's pod logs: `kubectl logs -f <pod-name>`
