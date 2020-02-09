FROM golang:1.13-alpine

RUN apk update && apk upgrade && apk add --no-cache bash git openssh curl

WORKDIR /go-gql-server

COPY . /go-gql-server/
RUN go mod download
# workarround for SP-291. See https://github.com/oxequa/realize/issues/253
RUN go get gopkg.in/urfave/cli.v2@master
RUN go get github.com/tockins/realize

CMD ["./scripts/run-dev.sh"]
