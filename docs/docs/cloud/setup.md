---
sidebar_position: 2
---

# Setup Your Workspace For Cloud

## Getting your cloud token

First, go to your workspace directory and run `mify cloud init`. If it is the
first time you ran this command, it will ask to you visit https://cloud.mify.io
to receive your token. Follow through the link and Sign in the Mify Cloud,
after that you'll be notified by email when you'll be able to access the cloud
console.

When you'll get the access, at first you'll be prompted to create your organization.
If you have no organization for your team, create it from here. If you want to
join existing organization, ping us on Slack.

After creating organization, you'll see dashboard with the next steps for
getting the token: ![](/img/docs/cloud-get-token.png)

Click on the "Generate New Service Token" button and copy result into the `mify
cloud init` prompt.

`mify cloud init` will link your project and organization to your workspace
and it'll update your `~/.kube/config` with credentials to connect to Mify Cloud's
Kubernetes.
