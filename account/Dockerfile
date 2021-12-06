FROM golang:1.17.2-buster AS build

RUN go version
ENV GOPATH=/
WORKDIR /src/
COPY ./ /src/

RUN go mod download; go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o /app ./cmd/main.go


FROM amd64/alpine:3

RUN apk update && apk upgrade \
    && apk add sqlite && apk add socat \
    && apk add --no-cache musl-dev gcc build-base \
    && apk add bash \
    && apk add --no-cache curl && apk add lsof

COPY --from=build /app /app
COPY ./.env /.env

WORKDIR /

RUN chmod +x app
RUN ls -la && pwd

CMD ["./app"]

HEALTHCHECK --interval=5s --timeout=3s --start-period=1s CMD curl --fail http://127.0.0.1:8001/health || exit 1
