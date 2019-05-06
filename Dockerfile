FROM golang:alpine
WORKDIR /app

COPY art ./art

RUN apk update && apk add --no-cache git
ADD animaniacs.go .
RUN go get -d -v ./...
RUN go build -o animaniacs .

ENTRYPOINT [ "/app/animaniacs" ]
