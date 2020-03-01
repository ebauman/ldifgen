FROM golang:1.13.4-alpine3.10 AS build
WORKDIR /go/src/github.com/ebauman/ldifgen
COPY . .

RUN go build

FROM alpine:latest
COPY --from=build /go/src/github.com/ebauman/ldifgen/ldifgen /bin/ldifgen

CMD /bin/ldifgen