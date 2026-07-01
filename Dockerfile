# Build stage
FROM golang:1.26-alpine3.20 AS builder

RUN apk add --no-cache make git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod verify && \
    make docker && \
    mv ./bin/nali-docker /build/nali

# Runtime stage - minimal
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /build/nali /usr/local/bin/nali
ENTRYPOINT ["nali"]