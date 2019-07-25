FROM golang:1.12.6 as builder

WORKDIR /go-modules

COPY . ./

# Building using -mod=vendor, which will utilize the v
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o vortex

FROM alpine:3.8

WORKDIR /root/

COPY --from=builder /go-modules/vortex .

ENTRYPOINT ["./vortex"]
