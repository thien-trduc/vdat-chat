# Build image
FROM golang:1.14-alpine as build
WORKDIR /go/src/app
COPY go.* ./
RUN go mod download
COPY . .

# Build Go Server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
          -ldflags='-w -s -extldflags "-static"' -a \
          -o /go/bin/chatserver ./cmd/chatserver

# Build Angular Webapp
FROM node:12-alpine AS angular-build
WORKDIR /usr/src/app
COPY ./website/package.json ./
RUN npm install
COPY ./website .
RUN npm run build:prod

# Target image
FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=build /go/bin/chatserver ./
COPY --from=build /go/src/app/migration/ ./migration/
COPY --from=angular-build /usr/src/app/dist ./public

EXPOSE 5000

CMD ["./chatserver"]
