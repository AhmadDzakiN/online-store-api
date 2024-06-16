FROM golang:1.22-alpine as Builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

# To configure your env value, please check your params/.env file and edit it first before creating the app's docker container
# If you want to change some of the env values, you need to build the image again

RUN go mod tidy

RUN go build -o main ./cmd

FROM alpine:latest

COPY --from=builder /app/main .

COPY params/.env .

# Change this config and the one in .env if you want to change the port
EXPOSE 1323

CMD ["./main"]