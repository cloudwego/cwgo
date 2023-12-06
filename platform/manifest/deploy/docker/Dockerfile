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

RUN go install github.com/cloudwego/cwgo@latest \
    && go install github.com/cloudwego/thriftgo@latest

RUN chmod 777 ./cwgo-platform

ENTRYPOINT ["./cwgo-platform"]
