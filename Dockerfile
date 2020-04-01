FROM golang:1.14-alpine as build
RUN apk add --no-cache git gcc musl-dev tzdata sqlite-dev
RUN go get github.com/gobuffalo/packr/v2/packr2
ADD . /app
WORKDIR /app
RUN packr2
RUN go test --tags "libsqlite3 linux"
RUN go build --tags "libsqlite3 linux"

FROM alpine:3.11
RUN apk add --no-cache tzdata sqlite-dev ca-certificates
RUN adduser -S -D -H -h /app kis3
COPY --from=build /app/kis3 /bin/
RUN mkdir /app && chown -R kis3 /app
USER kis3
WORKDIR /app
RUN mkdir data
VOLUME /app/data
EXPOSE 8080
CMD ["kis3"]
