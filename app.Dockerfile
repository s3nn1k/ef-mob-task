FROM golang:1.23-rc-alpine3.20 AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# build
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY docs ./docs
RUN go build -o ./bin/app ./cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY /.env /.env
COPY /migrations /migrations

CMD ["/app"]