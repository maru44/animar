FROM golang:1.17.1-alpine

# ENV ROOT=/go/src/app
# WORKDIR ${ROOT}

WORKDIR /deployments

COPY ../go.mod ../go.sum ./

RUN go mod download
COPY ../ .

EXPOSE 8080

CMD ["go", "run", "../cmd/project/main.go"]
