FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

ADD . ./
RUN go build

EXPOSE 8080

CMD ["go", "run", "main.go"]
