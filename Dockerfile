FROM golang:latest

WORKDIR /app

VOLUME [ "/src" ]

EXPOSE 8080

# Generate an SSL cert using Go's crypto/tls package
#   Satisfy requirement #1 - The login view must protect the security of data entry from eavesdroppers.
#   HTTPS encryption prevents MITM attacks
RUN go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

CMD go run src/main.go