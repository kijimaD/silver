###########
# builder #
###########

FROM golang:1.20-buster AS builder
RUN apt update \
    && apt install -y --no-install-recommends \
    upx-ucl sudo

WORKDIR /build
COPY . .

WORKDIR /build/example
RUN GO111MODULE=on CGO_ENABLED=0 go build -o /build/bin/silver \
    -ldflags='-w -s -extldflags "-static"' \
    . \
 && upx-ucl --best --ultra-brute /build/bin/silver

###########
# release #
###########

FROM gcr.io/distroless/static-debian11:latest AS release

COPY --from=builder /build/bin/silver /bin/
WORKDIR /workdir
ENTRYPOINT ["/bin/silver"]
