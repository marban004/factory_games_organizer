FROM golang:1.24

WORKDIR /app

ADD . /app
COPY crud_microservice/go.mod calculator_microservice/go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/crud_microservice
RUN go build -o main .
ENTRYPOINT [ "/app/crud_microservice/main" ]