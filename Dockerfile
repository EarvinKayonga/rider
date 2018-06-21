FROM        golang:1.10.3-alpine as builder
RUN         apk add -u --no-cache build-base git
ADD         .   /go/src/github.com/EarvinKayonga/rider
WORKDIR     /go/src/github.com/EarvinKayonga/rider
RUN         make build
RUN         strip /go/src/github.com/EarvinKayonga/rider/bin/*

FROM        alpine:latest as trip
EXPOSE      8082
WORKDIR     /root/
COPY        --from=builder /go/src/github.com/EarvinKayonga/rider/bin/trip rider
COPY        --from=builder /go/src/github.com/EarvinKayonga/rider/configuration.trip.yml configuration.yml
CMD         ["./rider", "--configuration", "configuration.yml"]

FROM        alpine:latest as gateway
EXPOSE      8080
WORKDIR     /root/
COPY        --from=builder /go/src/github.com/EarvinKayonga/rider/bin/gateway rider
COPY        --from=builder /go/src/github.com/EarvinKayonga/rider/configuration.gateway.yml configuration.yml
CMD         ["./rider", "--configuration", "configuration.yml"]

FROM        alpine:latest as bike
EXPOSE      8081
WORKDIR     /root/
COPY        --from=builder /go/src/github.com/EarvinKayonga/rider/bin/bike rider
COPY        --from=builder /go/src/github.com/EarvinKayonga/rider/configuration.bike.yml configuration.yml
CMD         ["./rider", "--configuration", "configuration.yml"]
