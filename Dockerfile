FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP=server

RUN go build -o /out/app ./cmd/${APP}


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /out/app /app/app

EXPOSE 8080

CMD ["/app/app"]