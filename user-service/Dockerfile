FROM golang:1.19.2-alpine as base
WORKDIR /build
COPY ./user-service/go.* ./
RUN go mod download
COPY ./user-service ./

FROM base as build
ENV OS linux
ENV ARCH amd64
COPY --from=base /go/pkg /go/pkg
COPY ./user-service /app
RUN CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -o /bin/app ./cmd/app

FROM scratch
COPY ./migrations/db /migrations
COPY --from=build /app/config /config
COPY --from=build /bin/app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
CMD ["/app"]
