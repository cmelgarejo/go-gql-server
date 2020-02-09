# Multistaged build production golang service
FROM golang:1.13-alpine as base

FROM base AS ci

RUN apk update && apk upgrade && apk add --no-cache git
RUN mkdir /build
ADD . /build/
WORKDIR /build

# Build prod
FROM ci AS build-env

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix \
    cgo -ldflags '-extldflags "-static"' -o server ./cmd/gql-server/

FROM alpine AS prod
RUN apk --no-cache add ca-certificates

COPY --from=build-env build/server ./gql-server/

CMD ["./gql-server/server"]