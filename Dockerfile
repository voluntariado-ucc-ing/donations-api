FROM golang:1.14

WORKDIR /go/src/github.com/voluntariado-ucc-ing/donations-api
COPY . .

RUN go get -d  -v ./...
RUN go install -v ./...

CMD ["donations-api"]

# Database Credentials
ENV DB_HOST=172.18.0.9

ENV DB_USER=postgres

ENV DB_PASS=ysl*gzzjic4Taok

ENV DB_NAME=voluntariado_ing

EXPOSE 8080
