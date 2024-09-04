# 打包依赖阶段使用golang作为基础镜像
FROM golang:latest as builder

LABEL Author="houyw<1327603193@qq.com>" Version="1.0" Description="build image for go-zero"

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/release

ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o hserver ./main.go

FROM scratch as prod
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /go/release/hserver /
COPY --from=builder /go/release/templates /templates
COPY --from=builder /go/release/conf /conf

ENV GIN_MODE=release \
    PORT=8080

EXPOSE 8080

CMD ["/hserver"]
