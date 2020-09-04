FROM golang:1.14

WORKDIR /go/src/github.com/voluntariado-ucc-ing/donations-api
COPY . .

RUN go get -d  -v ./...
RUN go install -v ./...

CMD ["donations-api"]

# Database Credentials

EXPOSE 8080
