FROM golang:1.16-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/sharesies-bot

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
COPY cmd/ ./cmd/

RUN go build -ldflags="-s -w -X main.Version=$VERSION" -o ./out/sharesies-bot ./cmd/ 


FROM alpine:3.9 

RUN apk add ca-certificates
RUN apk add --no-cache tzdata

COPY --from=build_base /tmp/sharesies-bot/out/sharesies-bot /app/sharesies-bot

CMD ["/app/sharesies-bot"]