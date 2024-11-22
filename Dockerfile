FROM golang:1.21.4-bookworm AS builder

RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    git \
    ca-certificates

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o build/ttv-bot .


FROM debian:bookworm-slim AS runner
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/build/ttv-bot .

ENTRYPOINT [ "./ttv-bot" ]






