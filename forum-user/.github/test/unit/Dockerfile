FROM golang:1.23

WORKDIR /root
COPY .. .
RUN go mod tidy

CMD go test -v ./test/unit/...