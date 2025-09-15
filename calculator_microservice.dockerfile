FROM golang:latest

WORKDIR /app

ADD ca.crt /etc/ssl/certs/
ADD calculator_microservice/application calculator_microservice/application
ADD calculator_microservice/handler calculator_microservice/handler
ADD calculator_microservice/microservice_logic_calculator calculator_microservice/microservice_logic_calculator
ADD calculator_microservice/docs calculator_microservice/docs
ADD calculator_microservice/custom_middleware calculator_microservice/custom_middleware
COPY calculator_microservice/go.mod calculator_microservice/go.sum calculator_microservice/main.go calculator_microservice/calculator_microservice_secret.pem calculator_microservice/calculator_microservice_cert.crt calculator_microservice/
WORKDIR /app/calculator_microservice
RUN go mod download
RUN go build -o main .
ENTRYPOINT [ "/app/calculator_microservice/main" ]