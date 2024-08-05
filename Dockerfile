FROM golang:1.21.11-alpine AS BASEIMAGE

WORKDIR /node/bot-middleware

COPY go.mod ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN go mod tidy
RUN swag init
RUN go build -o main .

FROM alpine
COPY --from=BASEIMAGE /node/bot-middleware ./
RUN apk add -U tzdata
ENV TZ Asia/Jakarta
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

EXPOSE 8181

CMD ["./main"]