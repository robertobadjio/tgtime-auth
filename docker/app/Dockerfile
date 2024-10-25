FROM golang:1.23 AS build

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../../cmd/auth/auth.go ./
COPY ../../pkg pkg
COPY ../../internal internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /tgtime-auth

FROM scratch

COPY --from=build /tgtime-auth .
COPY --from=ghcr.io/tarampampam/curl:8.6.0 /bin/curl /bin/curl

ENV TZ Europe/Moscow

CMD ["/tgtime-auth"]