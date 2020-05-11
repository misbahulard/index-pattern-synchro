FROM golang:latest as builder

WORKDIR /src
COPY . .

RUN go mod download
RUN go mod verify
RUN ls -alh
RUN CGO_ENABLED=0 go build -tags netgo -a -v -o synchro cmd/synchro/synchro.go

FROM alpine:3.8

WORKDIR /app
COPY --from=builder /src/synchro .

CMD [ "./synchro" ]
