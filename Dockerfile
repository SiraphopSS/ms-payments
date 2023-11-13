# Use an official Go runtime as a parent image
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./main.go

EXPOSE 8080

CMD [ "/main.go" ]