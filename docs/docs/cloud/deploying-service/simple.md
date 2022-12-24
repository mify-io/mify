---
sidebar_position: 0
---

# Simple deployment

For a quick deploy just run `mify cloud deploy` command in your workspace and
it will deploy every service and frontend in it.

If you want to deploy only one service, use:
```
$ mify cloud deploy <service-name>
```

To deploy in specific environment pass `-e` flag: (by default the command will deploy in `stage` environment)
```
$ mify cloud deploy <service-name> -e prod
```

Your frontend and backend services will be accessible via the corresponding
`https://<service-name>.app.mify.io` address. In the next page we'll describe
additional configuration for your services.

