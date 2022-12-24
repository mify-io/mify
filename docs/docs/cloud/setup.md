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

When you'll get the access, you'll see dashboard with the next steps for
getting the token: ![](/img/docs/cloud-get-token.png)

Click on the "Generate New Service Token" button and copy result into the `mify
cloud init` prompt.

## Registering workspace

After getting your token the tool will ask you for your unique project name, you can
either use your workspace directory name or choose something else. Then, select
environment, either `stage` or `prod`. When you get the message that you're
successfully registered your project you are ready to deploy your services.

