FROM docker.uclv.cu/golang:1.19.2 as dev

ENV PORT=2020

WORKDIR /go/src/

VOLUME [ "/go/src/data/bind/conf", "/go/src/data/bind/lib",  "/go/src/data/api" ]

EXPOSE 2020

ENTRYPOINT [ "./api-entrypoint.sh" ]

CMD ["go", "run", "."]
