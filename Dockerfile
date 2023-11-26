# FROM golang:1-alpine
# RUN apk add --no-cache git
# WORKDIR /src/app
# COPY . .
# RUN go build -o /server .
# ENTRYPOINT ["/server"]


####################
# Build Stage
FROM golang:1-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src/app
COPY . .
RUN go build -o server .

# Final Stage
FROM gcr.io/distroless/base
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /src/app/server /app/server
ENTRYPOINT ["/app/server"]

