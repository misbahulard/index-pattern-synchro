# Index Pattern Synchro - シンクロ

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

## Introduction

Index Pattern synchro is tool for synchronize elasticsearch indices with Kibana index pattern automatically. 

This tool will be the best for EFK implementation that you have many index and need you are too lazy to separate index pattern in Kibana. For example:

- log-nginx-2020.05.01
- log-nginx-2020.05.02
- log-nginx-2020.05.03
- log-tomcat-2020.05.01
- log-tomcat-2020.05.02
- log-app-2020.05.01

it will be generated an index pattern each of them so it will be:

- log-nginx*
- log-tomcat*
- log-app*

Of course this tool support for multiple **Kibana Spaces** or **Opendistro Tenants**!

## Configuration

This app need `config.yaml` to run, you can generate the config using the command line tool.

### Initialize or Generate the config file

To create `config.yaml` in current directory.

```
index-pattern-synchro init
```

To create `config.yaml` in specific directory, use the flag `-o` or `--output`.

```
index-pattern-synchro init -o /opt/index-pattern-synchro/config.yaml
```

*Configuration example is also in the `example/config.yaml`.*

### How config file works

This tool will discover the config file in these directories, the priority is from top to bottom.

- `/etc/index-pattern-synchro/config.yaml`
- `/opt/index-pattern-synchro/config.yaml`
- `$HOME/index-pattern-synchro/config.yaml` (user home directory)
- `./config.yaml` (current directory)

An explanation of each field of the config can be read in [config markdown](CONFIG.md).

## Build

### Build Binary

```
go build
```

### Docker Build image

```
docker built -f build/Dockerfile -t index-pattern-synchro:latest .
```

### Docker Run Container

```
docker run --rm --name index-pattern-synchro --mount type=bind,source="$(pwd)"/config.yaml,target=/app/config.yaml index-pattern-synchro
```

*Prebuild docker images exist in my docker hub [here](https://hub.docker.com/r/misbahulard/index-pattern-synchro).*

## Run

For running this app, you just need to execute the binary or run the container.

```
index-pattern-synchro run
```

You can also view the command line help, by adding the flag on it.

```
index-pattern-synchro -h
index-pattern-synchro run -h
```

Tips:
- Run the app with auto restart, so if the app exited it will be auto recovered.