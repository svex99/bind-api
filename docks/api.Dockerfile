FROM docker.uclv.cu/golang:1.19

ENV GO111MODULE=on
ENV GIN_MODE=release
ENV PORT=2020

WORKDIR /go/src/

COPY go.mod go.sum ./
COPY vendor vendor/

COPY . .

RUN go build -mod=vendor -o build/

VOLUME [ "/data" ]

EXPOSE 2020

CMD ["./build/bind-api"]
