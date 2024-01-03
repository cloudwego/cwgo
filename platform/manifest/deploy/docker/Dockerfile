FROM node:21.4.0-alpine3.18 as build-web

RUN corepack enable

COPY ./web /build
WORKDIR /build

RUN pnpm i \
    && pnpm run build

FROM golang:1.20.11-alpine3.18 as build-cmd

ARG GOPROXY
ENV GOPROXY=${GOPROXY}

COPY ./server /build
WORKDIR /build

ENV CGO_ENABLED=0

RUN go build -ldflags="-s -w" -o cwgo-platform cmd/command.go

FROM golang:1.20.11-alpine3.18 as run

ARG GOPROXY
ARG TZ

ENV GOPROXY=${GOPROXY}
ENV TZ=${TZ}

COPY --from=build-web /build/dist /app/dist
COPY --from=build-cmd /build/cwgo-platform /app/

WORKDIR /app

# add timezone file
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

# install protoc
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-linux-x86_64.zip \
    && unzip protoc-25.1-linux-x86_64.zip \
    && cp bin/protoc /usr/local/bin/protoc \
    && mkdir -p /usr/local/include/google \
    && cp -r include/google /usr/local/include/google

# install go clis
RUN go install github.com/cloudwego/cwgo@latest \
    && go install github.com/cloudwego/thriftgo@latest

RUN chmod 777 ./cwgo-platform

ENTRYPOINT ["./cwgo-platform"]
