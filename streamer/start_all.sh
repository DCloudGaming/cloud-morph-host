#!/usr/bin/env bash
#TODO: Async wait at each step
# We build this base image only once, shared among all containers
# For each RegisterApps request, different app_paths will be copied to this base image later.
# Copy-on-write https://stackoverflow.com/questions/36213646/why-are-containers-size-and-images-size-equivalent
# so spawn containers won't take much space initially regardless of number of containers
docker build -t syncwine .
go run main.go
# Start GUI App here
