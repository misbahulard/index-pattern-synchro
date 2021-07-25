FROM golang:latest as builder

WORKDIR /src
COPY . .

RUN go mod download
RUN go mod verify
RUN ls -alh
RUN CGO_ENABLED=0 go build -tags netgo -a -v -o index-pattern-synchro

FROM alpine:3.8

WORKDIR /app
COPY --from=builder /src/index-pattern-synchro .

CMD [ "./index-pattern-synchro" ]
