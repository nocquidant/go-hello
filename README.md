# go-hello

A tiny hello world HTTP application written in Go that is shipped as a docker container.

## API 

| Verb | Path   | Description |
| ---- | ------ | ----------- |
| GET  | /health | Display a 'UP' message |
| GET  | /hello | Display a 'hello' message |
| GET  | /request  | Call a remote service (get) using HELLO_REMOTE env variable |
| GET  | /info  | Display some technical informations |

# Docker

Docker images are published to: https://hub.docker.com/r/nocquidant/go-hello/