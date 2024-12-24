FROM golang:1.18-bullseye AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0
RUN go install github.com/beego/bee/v2@latest
ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor




ENV APP_HOME=/go/mathapp
RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"

ADD go.mod .
ADD go.sum .

COPY . .

EXPOSE 8010

RUN go mod vendor

CMD ["bee", "run"]