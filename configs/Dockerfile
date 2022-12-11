FROM golang:1.19 as build

WORKDIR /app

COPY . .

RUN go mod download && CGO_ENABLED=0 go build -o bin/rss-feed-bot ./cmd/rss-feed-bot


FROM alpine

WORKDIR /app

COPY --from=build /app/bin/rss-feed-bot .

ENTRYPOINT ["/app/rss-feed-bot"]
