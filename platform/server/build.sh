#!/usr/bin/env bash

api_files=(base idl registry repository template token)
cd cmd/api

for file in ${api_files[@]}; do
  hz update \
    --module github.com/cloudwego/cwgo/platform/server/cmd/api \
    --customize_package ../../template/api/package.yaml \
    --handler_by_method \
    --use ../../shared/kitex_gen/model \
    --model_dir internal/biz/model \
    --idl ../../idl/api/${file}.thrift \
    --handler_dir internal/biz/handler/${file} \
    --thriftgo "template=slim" ; \
done
