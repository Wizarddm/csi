#!/bin/bash

set -ex

# 编译
CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o build/nfs-csi cmd/main.go

image="nfs-csi:v1.0"

# 打包镜像
docker build -t $image . --platform linux/amd64

kind load docker-image $image

# 推送镜像
# docker push image