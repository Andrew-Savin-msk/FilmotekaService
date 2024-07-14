FROM ubuntu:24.04

RUN apt-get update && apt-get install -y wget tar curl gnupg lsb-release

ENV GOLANG_VERSION 1.22.2
RUN wget https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz -O go.tgz && \
    tar -C /usr/local -xzf go.tgz && \
    rm go.tgz

ENV PATH="/usr/local/go/bin:$PATH"

WORKDIR /go/restApiFilmoteka

ENV DOCKER_CONFIG_PATH=configs/server_docker.toml

COPY . .

RUN go mod tidy

CMD go run cmd/main.go
