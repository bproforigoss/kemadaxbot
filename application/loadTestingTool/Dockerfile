FROM golang:1.17
WORKDIR /gocode/src/github.com/bproforigoss/kemadaxbot/application/loadTestingTool
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a loadTestingTool.go

FROM scratch
EXPOSE 8080
COPY --from=0 /gocode/src/github.com/bproforigoss/kemadaxbot/application/loadTestingTool/loadTestingTool /loadTestingTool
COPY --from=0 /etc/ssl/certs/ca-certificates.crt  /etc/ssl/certs/
ENTRYPOINT ["/loadTestingTool"]
