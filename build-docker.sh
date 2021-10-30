#!/bin/bash -e
TAG=iamtio/mqtt2play:${VERSION:-1.0.0}
MACHINE=$(uname -m)
docker build -t $TAG .
docker image save -o docker_mqtt2play_$MACHINE.tar $TAG 