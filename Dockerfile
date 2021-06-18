FROM golang:1.16-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/auto-invest-sharesies

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
COPY cmd/ ./cmd/

RUN go build -ldflags="-s -w -X main.Version=$VERSION" -o ./out/auto-invest-sharesies ./cmd/ 


FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_base /tmp/auto-invest-sharesies/out/auto-invest-sharesies /app/auto-invest-sharesies

CMD ["/app/auto-invest-sharesies"]