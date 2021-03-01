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
# install cryptograhic dependency
RUN go get -v golang.org/x/crypto/bcrypt
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

ADD /Driver_Authentication /go/src/app/Driver_Authentication
ADD /Libraries /go/src/app/Libraries

# specifiy what directory we will use
WORKDIR /go/src/app/Driver_Authentication

# build
RUN go build -o driver_auth .

# entry point of service
CMD ["/go/src/app/Driver_Authentication/driver_auth"]


#docker run -p 3000:8080 -d v1
#docker build -t v1 -f Driver_Authentication/Dockerfile .