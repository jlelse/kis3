FROM golang:1.12-alpine as packr
RUN apk add --no-cache git
RUN go get github.com/gobuffalo/packr/v2/packr2

FROM golang:1.12-alpine as build
COPY --from=packr /go/bin/packr2 /go/bin
ADD . /app
WORKDIR /app
RUN apk add --no-cache git gcc musl-dev
RUN GO111MODULE=on packr2
RUN go build kis3.dev/kis3

FROM alpine:3.9
RUN adduser -S -D -H -h /app kis3
COPY --from=build /app/kis3 /app/
RUN chown -R kis3 /app
USER kis3
WORKDIR /app
RUN mkdir data
VOLUME ["/app/data"]
EXPOSE 8080
CMD ["./kis3"]
