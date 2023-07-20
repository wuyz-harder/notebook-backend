FROM golang:1.20
LABEL authors="wuyongzhen"

RUN mkdir -p "/data/app"
WORKDIR /data/app
# 认证文件也要复制进行
COPY ./cmd/cert.pem /data/app/cert.pem
COPY ./cmd/key.pem /data/app/key.pem
COPY ./cmd/main /data/app/main
ENTRYPOINT ["./main"]

