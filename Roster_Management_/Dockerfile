## We specify the base image we need for our
## go application
FROM golang:1.12.0-alpine3.9

# install gcc and git so we can install our dependencies
RUN apk add --update gcc musl-dev
RUN apk add git

# install sql dependency
RUN go get -v github.com/mattn/go-sqlite3
# install json web token dependency
RUN go get -v github.com/dgrijalva/jwt-go
# install http routing dependency
RUN go get -v github.com/gorilla/mux
# install http context dependency
RUN go get -v github.com/gorilla/context

# json form validation dependencies

RUN go get -v github.com/go-playground/locales/eu
RUN go get -v github.com/go-playground/universal-translator
RUN go get -v gopkg.in/go-playground/validator.v9
RUN go get -v gopkg.in/go-playground/validator.v9/translations/en
RUN go mod download

## create an /app directory within our docker image

RUN mkdir /go/src/app

# add local files the are required by the service

ADD /Roster_Management_ /go/src/app/Roster_Management
ADD /Libraries /go/src/app/Libraries

# specifiy what directory we will use

WORKDIR /go/src/app/Roster_Management

# build

RUN go build -o roster_management .

# entry point of service

CMD ["/go/src/app/Roster_Management/roster_management"]

#docker run -p 3002:8082 -it v3
#docker build -t v3 -f Roster_Management_/Dockerfile .
#docker run -p 3002:8082 -d v3
