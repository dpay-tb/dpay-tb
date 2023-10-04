FROM alpine:latest as builder
RUN apk add --no-cache \
    bash \
    git \
    curl

WORKDIR /build
RUN git clone https://github.com/tigerbeetle/tigerbeetle
RUN ./tigerbeetle/bootstrap.sh

FROM alpine:latest AS main
COPY --from=builder /build/tigerbeetle/tigerbeetle /tigerbeetle

ENTRYPOINT ["/tigerbeetle"]