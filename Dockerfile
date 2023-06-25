FROM golang:alpine as builder
WORKDIR /short-link

COPY  . .
RUN apk update
RUN CGO_ENABLED=0 GOOS=linux go build -o main -ldflags="-s -w" -a -installsuffix cgo ./cmd/main.go

FROM alpine
WORKDIR /short-link
COPY --from=builder ./short-link .

EXPOSE 8080
ARG db
ENV db ${db}
CMD ./main --db=${db}

