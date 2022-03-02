FROM golang:1.17
WORKDIR /gocode/src/github.com/bproforigoss/kemadaxbot/application/primeGenerator
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a primeGenerator.go

FROM scratch
EXPOSE 8080
COPY --from=0 /gocode/src/github.com/bproforigoss/kemadaxbot/application/primeGenerator/primeGenerator /primeGenerator
COPY --from=0 /etc/ssl/certs/ca-certificates.crt  /etc/ssl/certs/
ENTRYPOINT ["/primeGenerator"]
