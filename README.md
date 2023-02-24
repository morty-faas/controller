# Morty Gateway

Welcome to the gateway of the project Morty.

## Prerequisites

* [Golang](https://go.dev/doc/install) (`>1.18`)
* [RIK](https://github.com/polyxia-org/rik) (`>0.1.0`)

## Getting started

* Download the latest binary of the gateway from the 
[release page](https://github.com/polyxia-org/morty-gateway/releases).
* The serve command will start the gateway, and listen on the port 8080 (can 
be configured).

```bash
./morty-gateway serve
```

Flags specific to the serve command:

| Flag              | Description                | Default                 |
|-------------------|----------------------------|-------------------------|
| `-p --port`       | The port to listen on      | `8080`                  |
| `-c --controller` | The RIK controller address | `http://localhost:5000` |


Global flags:

| Flag              | Description           | Default                 |
|-------------------|-----------------------|-------------------------|
| `-h --help`       | Show help             |                         |
| `--config`        | The config file       | `morty-gateway.yaml`    |
