FROM golang:1.24

WORKDIR /app

ADD . /app
COPY calculator_microservice/go.mod calculator_microservice/go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/calculator_microservice
RUN go build -o main .
ENTRYPOINT [ "/app/calculator_microservice/main" ]