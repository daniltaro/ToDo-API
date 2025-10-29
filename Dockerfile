# 1 Stage
FROM golang:alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/server/main.go 

# 2 Stage
FROM alpine

RUN apk add --no-cache git
RUN apk add --no-cache go 
RUN go install github.com/pressly/goose/v3/cmd/goose@latest && apk del go

ENV PATH=$PATH:/root/go/bin

WORKDIR /build
COPY --from=builder /build/main .
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/entrypoint.sh .
COPY --from=builder /build/.env .

RUN chmod +x entrypoint.sh
ENTRYPOINT [ "./entrypoint.sh" ]