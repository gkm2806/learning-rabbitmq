FROM golang:1.14-alpine as build

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY src/main.go /src 

RUN go build main.go

FROM alpine as runtime

COPY --from=build /src/main /app/publisher

CMD [ "/app/publisher" ]
