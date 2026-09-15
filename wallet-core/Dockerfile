FROM golang:1.26

WORKDIR /app

RUN apt-get update \
    && apt-get install -y librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /usr/local/bin/walletcore ./cmd/walletcore

EXPOSE 8080
CMD ["walletcore"]