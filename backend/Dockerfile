FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go2-gkes-tpc .

FROM alpine:edge
WORKDIR /app
COPY --from=build /app/go2-gkes-tpc .
COPY --from=build /app/config.example.yaml ./config.yaml

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT ["/app/go2-gkes-tpc"]