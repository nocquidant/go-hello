# go-hello

A tiny hello world go http application that is shipped as a docker container.

## API 

| Verb | Path   | Description |
| ---- | ------ | ----------- |
| GET  | /hello | Display a hello message |
| GET  | /back  | Call back service using BACK_PORT env variable |
| GET  | /info  | Display some technical informations |

# Docker

Docker images are published to: https://hub.docker.com/r/nocquidant/go-hello/