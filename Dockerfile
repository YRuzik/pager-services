FROM golang:alpine AS builder
WORKDIR /build
ENV CGO_ENABLED=0

ADD ./go.mod .
ADD ./go.sum .
RUN go mod download

COPY . .

WORKDIR /build/cmd/pager
RUN go build -o pager -ldflags "-s -w"
FROM alpine

EXPOSE 4001
EXPOSE 5001
COPY --from=builder /build/cmd/pager/pager /pager
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
ENTRYPOINT ["/pager"]
CMD []