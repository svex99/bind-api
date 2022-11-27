FROM docker.uclv.cu/golang:1.19.2 as base

ENV CGO_ENABLED=0

WORKDIR /go/src

COPY . .

RUN go build -mod=vendor -o build/bind-api


FROM scratch as prod

ENV GIN_MODE=release

COPY --from=base /go/src/build/bind-api /usr/bin/bind-api

VOLUME [ "/data/bind/conf", "/data/bind/lib", "/data/api" ]

EXPOSE 2020

CMD ["bind-api"]
