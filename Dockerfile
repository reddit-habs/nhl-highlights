FROM docker.io/golang:1.19-alpine3.17 AS builder
WORKDIR /build
RUN apk add --no-cache build-base
COPY . /build
RUN go mod tidy
RUN ./tools/build

FROM docker.io/alpine:3.17
COPY --from=builder /build/build/nhl-highlights /usr/local/bin/nhl-highlights
WORKDIR /data
EXPOSE 9999
ENTRYPOINT ["nhl-highlights"]
CMD ["serve", "--incremental"]
