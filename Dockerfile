FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/octyne ./cmd/octyne

FROM alpine:3.22

RUN adduser -D -H -u 10001 octyne

WORKDIR /app

COPY --from=build /out/octyne /usr/local/bin/octyne

USER octyne

EXPOSE 3000

ENTRYPOINT ["octyne"]