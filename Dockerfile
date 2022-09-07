FROM gcr.io/distroless/static-debian11 as prod
CMD ["/app"]

FROM golang:1.18 as build
WORKDIR /app

COPY go.mod .

RUN go mod download
RUN go mod verify

COPY . .
RUN go vet -v
RUN go test -v
RUN CGO_ENABLED=0 go build -o /app/bin/rosina

FROM prod
COPY --from=build /app/bin/rosina /app
