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

ADD /Driver_Allocation /go/src/app/Driver_Allocation
ADD /Libraries /go/src/app/Libraries

# specifiy what directory we will use
WORKDIR /go/src/app/Driver_Allocation

# build
RUN go build -o driver_allocation .

# entry point of service
CMD ["/go/src/app/Driver_Allocation/driver_allocation"]


#docker run -p 3003:8083 -d v4
#docker build -t v4 -f Driver_Allocation/Dockerfile .