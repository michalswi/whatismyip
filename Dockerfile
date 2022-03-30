ARG GOLANG_VERSION
ARG ALPINE_VERSION

# build
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

ARG APPNAME

RUN apk --no-cache add make

WORKDIR /app

COPY main.go main.go
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download
RUN make build

# execute
FROM alpine:${ALPINE_VERSION}

ARG APPNAME

ENV SERVER_PORT ""
ENV PPROF_PORT ""

RUN adduser -D -h /dummy dummy
USER dummy
WORKDIR /dummy

COPY --from=builder /app/${APPNAME} ./${APPNAME}

CMD ["./whatismyip2"]
