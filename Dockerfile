FROM golang:alpine
WORKDIR /app

COPY art ./art

RUN apk update && apk add --no-cache git
ADD animaniacs.go go.mod go.sum ./
RUN go build -o animaniacs .

ENTRYPOINT [ "/app/animaniacs" ]
