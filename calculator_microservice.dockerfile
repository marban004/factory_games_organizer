FROM golang:1.24

WORKDIR /app

ADD calculator_microservice/application calculator_microservice/application
ADD calculator_microservice/handler calculator_microservice/handler
ADD calculator_microservice/microservice_logic_calculator calculator_microservice/microservice_logic_calculator
COPY calculator_microservice/go.mod calculator_microservice/go.sum calculator_microservice/main.go calculator_microservice/
WORKDIR /app/calculator_microservice
RUN go mod download
RUN go build -o main .
ENTRYPOINT [ "/app/calculator_microservice/main" ]