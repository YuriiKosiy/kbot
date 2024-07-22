FROM quay.io/projectquay/golang:1.22 as builder

WORKDIR /go/src/app
COPY . .
RUN make build

FROM scratch
WORKDIR /app
COPY --from=builder /go/src/app/kbot .
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY config.json /app/config.json
ENTRYPOINT ["./kbot", "start"]