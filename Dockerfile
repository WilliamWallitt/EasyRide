FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
RUN go get github.com/WilliamWallitt/EasyRide
RUN cd /build && git clone https://github.com/WilliamWallitt/EasyRide

RUN cd /build/EasyRide && go build

EXPOSE 10000

ENTRYPOINT ["/build/EasyRide/main"]
