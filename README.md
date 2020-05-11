# Index Pattern Synchro - シンクロ
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

- log-nginx-*
- log-tomcat-*
- log-app-*

## Configuration
This app need `config.yaml` to run, the example config is in `config/synchro/config.yaml`.

```
elasticsearch: [elasticsearch host]
kibana: [kibana host]
indexPattern: [base index-pattern, ex: "log-"]
kibanaMaxPage: [kibana max pagination for get list, ex: 1000]
interval: [interval for synchronize in minutes, ex: 3]
```

## Run
make sure you have `config.yaml` in current path

```
go run cmd/synchro/synchro.go
```

## Install

```
go install cmd/synchro/synchro.go
```

## Docker Build image

```
docker built -f build/synchro/Dockerfile -t yourtag/index-pattern-synchro .
```

## Docker Run

```
docker run --rm --name synchro --mount type=bind,source="$(pwd)"/config.yaml,target=/app/config.yaml yourtag/index-pattern-synchro
```