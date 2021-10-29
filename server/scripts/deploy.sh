#!/usr/bin/env bash
## Install nginx and sqlite3
# sudo apt-get install nginx sqlite3

## Create db prod.db and table schema
# sqlite3 prod.db

## Change config.yaml path in main.go
## Change db path in env.go

## Build update
# go build

## Create service file to deploy
cp declo_backend.service /etc/systemd/system/declo_backend.service
cp declo_nginx.com /etc/nginx/sites-enabled/declo.co

sudo systemctl daemon-reload
sudo systemctl restart declo_backend
#sudo systemctl restart nginx
