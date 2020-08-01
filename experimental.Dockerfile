# syntax=docker/dockerfile:experimental
# build go
FROM golang:alpine as builder
ARG version=next
RUN apk update && apk add --no-cache git gcc g++ libc-dev musl-dev
RUN addgroup -S appgroup && adduser -S -D -H -h /app -G appgroup appuser
COPY . $GOPATH/src/github.com/go-simple-startpage
WORKDIR $GOPATH/src/github.com/go-simple-startpage

ENV GO111MODULE=on
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -ldflags='-linkmode external -extldflags "-static" -s -w -X main.version='${version}'' -a -installsuffix cgo -o /go/bin/go-simple-startpage .

# build angular
FROM node:10-alpine as frontend
WORKDIR /app
COPY ui/package*.json ./
RUN npm ci
COPY ./ui .
RUN npm run build-prod

# build final image
FROM scratch
EXPOSE 3000
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --chown=appuser:appgroup --from=builder /go/bin/go-simple-startpage /app/
USER appuser
WORKDIR /app
COPY ./config.yaml .
COPY --from=frontend /app/dist ./ui/dist
ENTRYPOINT ["/app/go-simple-startpage"]
