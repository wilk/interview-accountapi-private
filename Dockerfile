FROM golang:1.13.5

RUN apt-get update && apt-get install -y netcat-openbsd

WORKDIR /go/src/app
COPY . .

CMD ["go", "test"]
