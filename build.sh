#!/bin/bash

set -ex

# 编译
CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o build/toy-lustre-csi cmd/main.go

image="toy-lustre-csi:1.0"
remoteImage=172.31.4.89/test/$image

# 打包镜像
docker build -t $image . --platform linux/amd64

# 推送镜像
docker tag $image $remoteImage
docker push $remoteImage
# docker push image