FROM golang:1.16
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o entrypoint ./cmd/payment-system
EXPOSE 8080
CMD ["/app/entrypoint"]
