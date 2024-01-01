#!/usr/bin/env bash

api_files=(base idl registry repository template token)

function api_gen() {
  cd cmd/api
  for file in ${api_files[@]}; do
    hz update \
        --module github.com/cloudwego/cwgo/platform/server/cmd/api \
        --customize_package ../../template/api/package.yaml \
        --handler_by_method \
        --use ../../shared/kitex_gen \
        --idl ../../idl/api/${file}.thrift \
        --handler_dir internal/biz/handler/${file} \
        --model_dir=../../shared/kitex_gen \
        --thriftgo "template=slim"
  done

}

function reuse_kx_model() {
  cd cmd/api
  for file in ${api_files[@]}; do
    hz model -module github.com/cloudwego/cwgo/platform/server/cmd/api \
     	--model_dir=../../shared/kitex_gen \
     	-t=ignore_initialisms \
     	-t=gen_setter \
     	-t=gen_deep_equal \
     	-t=compatible_names \
     	-t=frugal_tag \
     	-t='template=slim' \
     	--idl ../../idl/api/${file}.thrift
  done
}

if [ "$1" == "api_gen" ]; then
    api_gen
elif [ "$1" == "reuse_kx_model" ]; then
    reuse_kx_model
else
    echo "Invalid option. Usage: $0 {api_gen|reuse_kx_model}"
fi
