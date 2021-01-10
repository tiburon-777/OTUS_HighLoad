FROM golang:1.14 as builder
RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN go get -d ./cmd/.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app ./cmd/.

FROM alpine:latest
ENV APP_SERVER_PORT $APP_SERVER_PORT
WORKDIR /
COPY --from=builder /app/app ./sbin
COPY --from=builder /app/templates ./templates
EXPOSE ${APP_SERVER_PORT}
ENTRYPOINT ["app"]
