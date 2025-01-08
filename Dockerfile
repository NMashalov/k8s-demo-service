FROM golang:1.23-alpine
WORKDIR /build
COPY ./demo /build/
CMD ["./demo"]