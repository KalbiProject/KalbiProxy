FROM golang:alpine

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o cmd/kalbi-proxy/main.go .
WORKDIR /dist
RUN cp /build/main .
EXPOSE 5060
CMD ["/dist/main"]