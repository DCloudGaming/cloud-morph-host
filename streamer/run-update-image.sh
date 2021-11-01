#!/usr/bin/env bash
mkdir apps/$2
cp -r $1 /apps/$2
docker build -t syncwine --build-arg AppName=$2 -f Dockerfile.update .

