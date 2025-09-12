FROM golang:latest

WORKDIR /app

ADD crud_microservice/application crud_microservice/application
ADD crud_microservice/handler crud_microservice/handler
ADD crud_microservice/microservice_logic_crud crud_microservice/microservice_logic_crud
ADD crud_microservice/docs crud_microservice/docs
ADD crud_microservice/custom_middleware crud_microservice/custom_middleware
COPY crud_microservice/go.mod crud_microservice/go.sum crud_microservice/main.go crud_microservice/crud_microservice_secret.pem crud_microservice/crud_microservice_cert.crt crud_microservice/
WORKDIR /app/crud_microservice
RUN go mod download
RUN go build -o main .
ENTRYPOINT [ "/app/crud_microservice/main" ]