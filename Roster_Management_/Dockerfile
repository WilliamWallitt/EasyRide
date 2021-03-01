## We specify the base image we need for our
## go application
FROM golang:1.12.0-alpine3.9


#maybe

RUN apk add --update gcc musl-dev
RUN apk add git

RUN go get -v github.com/mattn/go-sqlite3
RUN go get -v github.com/dgrijalva/jwt-go
RUN go get -v github.com/gorilla/mux
RUN go get -v github.com/gorilla/context


RUN go get -v github.com/go-playground/locales/eu
RUN go get -v github.com/go-playground/universal-translator
RUN go get -v gopkg.in/go-playground/validator.v9
RUN go get -v gopkg.in/go-playground/validator.v9/translations/en
RUN go mod download
## We create an /app directory within our
## image that will hold our application source

RUN mkdir /go/src/app

ADD /Roster_Management_ /go/src/app/Roster_Management
ADD /Libraries /go/src/app/Libraries

WORKDIR /go/src/app/Roster_Management

RUN go build -o roster_management .

CMD ["/go/src/app/Roster_Management/roster_management"]


#docker build -t v3 -f Roster_Management_/Dockerfile .
#docker run -p 3002:8082 -d v3
