FROM golang:1.15.6-alpine as builder

ENV GO111MODULE "on"

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . /build
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build  -o app ./go-services/historical-events/cmd

FROM alpine:3.12.3 as base
COPY --from=builder /build/app /historical-events
COPY --from=builder /build/go-services/historical-events/api/proto /go-services/historical-events/api/proto
COPY --from=builder /build/go-services/historical-events/api/swagger-ui /go-services/historical-events/api/swagger-ui

EXPOSE 8080 9999
ENTRYPOINT [ "/historical-events" ]
