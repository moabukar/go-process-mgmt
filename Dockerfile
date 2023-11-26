FROM golang:1-alpine
RUN apk add --no-cache git
WORKDIR /src/app
COPY . .
RUN go build -o /server .
ENTRYPOINT ["/server"]
