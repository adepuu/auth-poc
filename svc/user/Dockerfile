FROM golang:1.18-alpine
RUN apk update && apk add --virtual build-dependencies build-base gcc wget git
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -gcflags="all=-N -l" -o main .
CMD ["/app/main --env development  --pid-file user.pid &>>user.log & disown"]