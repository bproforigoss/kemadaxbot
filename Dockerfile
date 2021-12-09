FROM golang:1.17
WORKDIR /gocode/src/github.com/bproforigoss/kemadaxbot
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a kemadaxbot.go

FROM scratch
EXPOSE 8080
COPY --from=0 /gocode/src/github.com/bproforigoss/kemadaxbot/kemadaxbot /kemadaxbot
COPY --from=0 /etc/ssl/certs/ca-certificates.crt  /etc/ssl/certs/
ENTRYPOINT ["/kemadaxbot"]
