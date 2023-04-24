# Morty Controller

This repository contains the source code of the Morty controller.

## Prerequisites

- [Golang](https://go.dev/doc/install) (`>1.19`)
- [RIK](https://github.com/polyxia-org/rik) (`>0.1.0`)

## Getting started

- Download the latest binary of the gateway from the
  [release page](https://github.com/morty-faas/controller/releases).
- The serve command will start the gateway, and listen on the port 8080 (can
  be configured).

```bash
./morty-gateway
```

By default, the controller will start on port `8080` and will target a RIK cluster on `http://localhost:5000`.

## Public API

This component provides its own OpenAPI specification available [here](./api/spec/openapi.yml). Based on this specification, client is generated and it is available to use in your projects to avoid the need to code the interaction with the Morty APIs. You can find more details about this client in [pkg/client](./pkg/client/).

In addition to all the routes listed in the specification, there is a dedicated route for invoking functions, `/functions/:name/invoke` where `:name` is the name of a function to invoke.
This route accepts **any HTTP methods** and will proxy the entire incoming request to an instance of the function directly. **This route will not be bundled in the autogenerated client nor listed in the specification.**

## Configuration

This component supports configuration over environments variables and YAML configuration file. By default at runtime, the component will try to retrieve the configuration from file `controller.yaml` present in the following directories :

- `/etc/morty`
- `$HOME/.morty`
- `.`

For development, you can create a configuration file `controller.yaml` at the root of the project and place the following content in it :

```yaml
port: 8080

orchestrator:
  rik:
    cluster: http://localhost:5000
# state:
#   redis:
#     addr: localhost:6379
```

> Note that `state` stanza is commented. By default, the `memory` engine will be loaded by the application. Add the required configuration for the adapter you want to use. For example here, if you uncomment, the application will try to initialize a `redis` engine adapter using the given `addr`.

If you wish to override a configuration through environment variables, for example `orchestrator.rik.cluster`, export the following environment variable :

```bash
export MORTY_CONTROLLER_ORCHESTRATOR_RIK_CLUSTER="mycustomaddr"
```

Finally, the default log level is `INFO`. If you want to update it, export the following environmnent variable with the level you want :

```bash
export MORTY_CONTROLLER_LOG=trace
```
