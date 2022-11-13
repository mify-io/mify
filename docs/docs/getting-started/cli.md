---
sidebar_position: 3
---

# Mify CLI Reference

Global flags:

`-h, --help`          help for mify.

`-p, --path string`   Path to workspace (if you're calling CLI outside directory).

`-v, --verbose`       Show verbose output.

Commands:

`mify init <project-name>` Initialize new workspace.

`mify add service <service-name> --language go|python` Generate service.

`mify add frontend <frontend-service-name>` Generate frontend.

`mify add client <service-name> --to <service-name>` Generate a client for service.

`mify remove client <service-name> --to <service-name>` Generate a client for service.

`mify generate <service-name>` Generate or regenerate code in workspace.
