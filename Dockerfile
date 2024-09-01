FROM quay.io/projectquay/golang:1.22 as builder

WORKDIR /go/src/app
COPY . .
RUN make build

FROM alpine:latest as tzdata
RUN apk --no-cache add tzdata

FROM scratch
WORKDIR /app
COPY --from=builder /go/src/app/kbot .
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=tzdata /usr/share/zoneinfo /usr/share/zoneinfo
ENTRYPOINT ["./kbot", "start"]
