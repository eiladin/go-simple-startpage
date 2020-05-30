# # build go
# FROM golang:alpine as builder
# RUN apk update && apk add --no-cache git gcc g++ libc-dev musl-dev
# RUN addgroup -S appgroup && adduser -S -D -H -h /app -G appgroup appuser
# COPY . $GOPATH/src/github.com/go-simple-startpage
# WORKDIR $GOPATH/src/github.com/go-simple-startpage

# ENV GO111MODULE=on
# RUN go mod download

# RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -ldflags='-linkmode external -extldflags "-static"' -a -installsuffix cgo -o /go/bin/go-simple-startpage .

# build angular
FROM node:12-alpine as frontend
WORKDIR /app
COPY ui/package*.json ./
RUN npm ci
COPY ./ui .
RUN npm run build -- --prod --aot --no-progress --build-optimizer

# build final image
# FROM scratch
# EXPOSE 3000
# COPY --from=builder /etc/passwd /etc/passwd
# COPY --from=builder /etc/group /etc/group
# COPY --chown=appuser:appgroup --from=builder /go/bin/go-simple-startpage /app/
# USER appuser
# WORKDIR /app
# COPY ./config.yaml .
# COPY --from=frontend /app/dist ./ui/dist
# ENV GIN_MODE=release
# ENTRYPOINT ["/app/go-simple-startpage"]
