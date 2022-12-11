# BUILD
FROM golang:1.19.4-alpine3.17 AS build-env

ENV APP_NAME mensa-restful
ENV CMD_PATH main.go

WORKDIR /$APP_NAME
COPY . .
RUN go build -o $APP_NAME

# RUN
FROM alpine:3.17.0 AS prod-env

ENV APP_NAME mensa-restful

COPY --from=build-env /$APP_NAME .
EXPOSE 8080
CMD ./$APP_NAME