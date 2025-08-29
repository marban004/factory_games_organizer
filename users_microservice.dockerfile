FROM golang:latest

WORKDIR /app

ADD users_microservice/application users_microservice/application
ADD users_microservice/handler users_microservice/handler
ADD users_microservice/microservice_logic_users users_microservice/microservice_logic_users
ADD users_microservice/docs users_microservice/docs
COPY users_microservice/go.mod users_microservice/go.sum users_microservice/main.go users_microservice/users_microservice_secret.pem users_microservice/users_microservice_cert.crt users_microservice/
WORKDIR /app/users_microservice
RUN go mod download
RUN go build -o main .
ENTRYPOINT [ "/app/users_microservice/main" ]