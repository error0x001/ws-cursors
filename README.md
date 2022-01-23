# ws-cursors

![Coverage](https://img.shields.io/badge/Coverage-76.8%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/error0x001/ws-cursors)](https://goreportcard.com/report/github.com/error0x001/ws-cursors)

![build](https://github.com/error0x001/ws-cursors/actions/workflows/build.yml/badge.svg)
![deploy](https://github.com/error0x001/ws-cursors/actions/workflows/deploy.yml/badge.svg)

## Description

A simple go app that can show users' cursors in real-time.

## Instance

The real service is available on Heroku by address - https://ws-cursors.herokuapp.com/

## ENVs

- `ADDRESS=0.0.0.0` - app's host
- `PORT=4567` - a port which uses by app
- `TEMPLATE_PATH=/go/bin/templates/index.html` - a path to the template
- `SHUTDOWN_TIME=5` - graceul shutdown timeout
- `IS_SSL_USING=0` - a sign if using SSL

## Run

You can run the app via docker, just execute `docker-compose up -d`. After that, open your browser `0.0.0.0:4567`
