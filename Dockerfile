FROM golang:1.14 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./src .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server /bin
EXPOSE 80
CMD ["server"]
