FROM golang:alpine3.12 as build
RUN mkdir -p /go/metaldata && \
        apk add upx
COPY ./ /go/metaldata
RUN cd /go/metaldata && \
        go mod vendor && \
        CGO_ENABLED=0 go build -x -ldflags '-s -w' -o /metaldata . && \
        ls -alh /metaldata && \
        upx /metaldata && \
        ls -alh /metaldata


FROM scratch
COPY --from=build /metaldata /
ENV MD_FSBASE=/data
ENV MD_BIND=:80
CMD ["/metaldata"]
