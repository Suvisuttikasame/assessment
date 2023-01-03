FROM golang:1.19-alpine as base-simple

WORKDIR /usr/src/app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v ./...

RUN go build -o ./out/go-app .

#section above is doing same thing as simple version

#section below, copy built file to alpine:3.16.2 and run built file
#this will save a lot of storage

FROM alpine:3.16.2

COPY --from=base-simple /usr/src/app/out/go-app /app/go-app

CMD ["/app/go-app"]


