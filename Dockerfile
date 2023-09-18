###########
# builder #
###########

FROM golang:1.20-buster AS builder
RUN apt update \
    && apt install -y --no-install-recommends \
    upx-ucl

WORKDIR /build
COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 go build -o ./bin/systemcheck \
    -ldflags='-w -s -extldflags "-static"' \
    . \
 && upx-ucl --best --ultra-brute ./bin/systemcheck

###########
# release #
###########

FROM gcr.io/distroless/static-debian11:latest AS release

COPY --from=builder /build/bin/systemcheck /bin/
WORKDIR /workdir
ENTRYPOINT ["/bin/systemcheck"]
