#!/bin/bash

echo "Setting docker repo..."
docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD

echo "Setting up build variables..."
GIT_COMMIT=$(git rev-list -1 HEAD)

echo "Setting up buildx..."
docker buildx version
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx create --name multiarch --use

echo "Building..."
docker buildx build \
--platform linux/amd64,linux/arm64,linux/arm/v7 \
-t supporttools/probe-a-node:${DRONE_BUILD_NUMBER} \
-t supporttools/probe-a-node:latest \
--cache-from supporttools/probe-a-node:latest \
--build-arg GIT_COMMIT=${DRONE_COMMIT} \
--build-arg GIT_BRANCH=${DRONE_BRANCH} \
--push -f Dockerfile .
