FROM golang:1.15-alpine AS build-env
RUN apk add ca-certificates upx
ADD . /goddns/
RUN cd /goddns \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o goddns goddns.go \
    && upx -9 goddns

FROM scratch
COPY --from=build-env /goddns/goddns /
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/goddns"]
