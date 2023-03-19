FROM golang:latest AS builder

ARG go_token

ENV GOPRIVATE="github.com/cocov-ci"
ENV CGO_ENABLED=0

RUN adduser --home /cocov cocov

RUN groupadd -f cocov
RUN usermod -a -G cocov cocov
RUN mkdir /cocov/app
WORKDIR /cocov/app

RUN /bin/bash

