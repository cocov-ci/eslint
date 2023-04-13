FROM golang:latest AS builder

ARG go_token

RUN apt-get update
RUN apt-get -y install openssl
RUN git config --global url."https://oauth2:$go_token@github.com/".insteadOf "https://github.com/"

ENV GOPRIVATE="github.com/cocov-ci"
ENV CGO_ENABLED=0

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build cmd/main.go

FROM debian:latest

COPY --from=builder /app/main /bin/eslint

RUN adduser --home /cocov cocov

RUN groupadd -f cocov
RUN usermod -a -G cocov cocov

USER cocov

CMD ["eslint"]
