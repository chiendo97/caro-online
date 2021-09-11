FROM golang:1.17-alpine3.13 as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git make \
    util-linux pciutils usbutils coreutils binutils findutils grep \
    build-base gcc abuild binutils binutils-doc gcc-doc

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make server

# FROM alpine:3.13

# RUN apk update && apk upgrade && \
#     apk add --no-cache git make

# WORKDIR /app

# EXPOSE 9091

# COPY --from=builder /app .

CMD ["make", "run_server"]
