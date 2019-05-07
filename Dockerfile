FROM alpine:edge as build-base
RUN apk add --no-cache go git gcc musl-dev sqlite

FROM build-base as packr
RUN go get github.com/gobuffalo/packr/v2/packr2

FROM build-base as build
COPY --from=packr /root/go/bin/packr2 /bin/
ADD . /app
WORKDIR /app
RUN GO111MODULE=on packr2
RUN go test
RUN go build --tags "libsqlite3 linux" github.com/kis3/kis3

FROM alpine:3.9
RUN apk add --no-cache tzdata ca-certificates sqlite
RUN adduser -S -D -H -h /app kis3
COPY --from=build /app/kis3 /bin/
RUN mkdir /app && chown -R kis3 /app
USER kis3
WORKDIR /app
RUN mkdir data
VOLUME ["/app/data"]
EXPOSE 8080
CMD ["kis3"]
