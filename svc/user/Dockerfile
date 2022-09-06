FROM golang:1.18-alpine
RUN apk update && apk add --virtual build-dependencies build-base gcc wget git
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod tidy
RUN go build -gcflags="all=-N -l" -o main .
EXPOSE 8002 30002
CMD ["/app/main", "--env", "development"]