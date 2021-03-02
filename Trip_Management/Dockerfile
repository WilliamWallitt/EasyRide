## We specify the base image we need for our
## go application
FROM golang:1.12.0-alpine3.9

# install gcc and git so we can install our dependencies
RUN apk add --update gcc musl-dev
RUN apk add git

# install http routing dependency
RUN go get -v github.com/gorilla/mux

# json form validation dependencies
RUN go get -v github.com/go-playground/locales/eu
RUN go get -v github.com/go-playground/universal-translator
RUN go get -v gopkg.in/go-playground/validator.v9
RUN go get -v gopkg.in/go-playground/validator.v9/translations/en
RUN go mod download

## create an /app directory within our docker image

RUN mkdir /go/src/app

# add local files the are required by the service

ADD /Trip_Management /go/src/app/Trip_Management
ADD /Libraries /go/src/app/Libraries

# specifiy what directory we will use
WORKDIR /go/src/app/Trip_Management

# build
RUN go build -o trip_management .

# entry point of service
CMD ["/go/src/app/Trip_Management/trip_management"]


#docker run -p 3003:8083 -d v4
#docker build -t v4 -f Trip_Management/Dockerfile .