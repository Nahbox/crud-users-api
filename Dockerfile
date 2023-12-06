FROM golang:1.21-alpine3.18 as builder

COPY . /go/src/github.com/Nahbox/crud-users-api

WORKDIR /go/src/github.com/Nahbox/crud-users-api
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/bin/crud-users-api github.com/Nahbox/crud-users-api/cmd

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/Nahbox/crud-users-api/ /go/src/github.com/Nahbox/crud-users-api/
COPY --from=builder /usr/bin/crud-users-api /usr/bin/crud-users-api
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/crud-users-api"]