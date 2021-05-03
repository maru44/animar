FROM golang:1.16.3-alpine
RUN apk update apk add git

ENV ROOT=/go/src/app
WORKDIR ${ROOT}

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]
