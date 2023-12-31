rpc_gen:
	kitex \
	  -frugal-pretouch \
	  -thrift frugal_tag,template=slim \
      -module github.com/cloudwego/cwgo/platform/server \
      -gen-path shared/kitex_gen \
      idl/agent/service.thrift ; \
    cd cmd/agent/ && \
    kitex -module github.com/cloudwego/cwgo/platform/server \
      -frugal-pretouch \
      -use github.com/cloudwego/cwgo/platform/server/shared/kitex_gen \
      -template-dir ../../template/agent \
      ../../idl/agent/service.thrift

FILE_NAME ?= ""

api_gen:
	bash build.sh

