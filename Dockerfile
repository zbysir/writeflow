FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags \"-static\"" -o writeflow main.go
RUN chmod +x writeflow

FROM registry.cn-hangzhou.aliyuncs.com/bysir/alpine-shanghai:latest

COPY --from=builder /go/src/app/writeflow /

ENTRYPOINT ["/writeflow"]