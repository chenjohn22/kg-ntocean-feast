FROM golang:1.24-alpine AS build

WORKDIR /src/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend ./
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server \
 && CGO_ENABLED=0 go build -o /out/migrate ./cmd/migrate

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Taipei
WORKDIR /app/backend
COPY --from=build /out/server /usr/local/bin/server
COPY --from=build /out/migrate /usr/local/bin/migrate
COPY backend/migrations ./migrations
COPY backend/data /app/backend/data
COPY random_codes_60000.json /app/random_codes_60000.json
EXPOSE 8081
CMD ["/usr/local/bin/server"]
