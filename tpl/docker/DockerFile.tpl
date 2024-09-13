FROM golang:{{.GoVersion}}alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
{{if .EnableGoProxy}}ENV GOPROXY {{.GoProxy}}
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
{{end}}{{if .HasTimezone}}
RUN apk update --no-cache && apk add --no-cache tzdata
{{end}}
WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY {{.GoRelPath}}/etc /app/etc
RUN go build -ldflags="-s -w" -o /app/{{.ExeFile}} {{.GoMainFrom}}


FROM {{.BaseImage}}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
{{if .HasTimezone}}COPY --from=builder /usr/share/zoneinfo/{{.Timezone}} /usr/share/zoneinfo/{{.Timezone}}
ENV TZ {{.Timezone}}
{{end}}
WORKDIR /app
COPY --from=builder /app/{{.ExeFile}} /app/{{.ExeFile}}{{if .Argument}}
COPY --from=builder /app/etc /app/etc{{end}}
{{if .HasPort}}
EXPOSE {{.Port}}
{{end}}
CMD ["./{{.ExeFile}}"{{.Argument}}]
