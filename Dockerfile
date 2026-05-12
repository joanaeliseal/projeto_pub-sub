FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /bin/broker ./broker
RUN go build -o /bin/loadbalancer ./loadbalancer
RUN go build -o /bin/publisher-cambio ./exemplos/publisher-cambio
RUN go build -o /bin/publisher-cripto ./exemplos/publisher-cripto
RUN go build -o /bin/subscriber-importadores ./exemplos/subscriber-importadores
RUN go build -o /bin/subscriber-investidores ./exemplos/subscriber-investidores

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/broker /bin/broker
COPY --from=builder /bin/loadbalancer /bin/loadbalancer
COPY --from=builder /bin/publisher-cambio /bin/publisher-cambio
COPY --from=builder /bin/publisher-cripto /bin/publisher-cripto
COPY --from=builder /bin/subscriber-importadores /bin/subscriber-importadores
COPY --from=builder /bin/subscriber-investidores /bin/subscriber-investidores
