FROM golang:latest AS build
WORKDIR /app
COPY .. /app/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o . cmd/sso/main.go

FROM scratch
COPY --from=build /app/main /app/config/demo.yaml /app/
ENV CONFIG_PATH="/app/demo.yaml"
ENTRYPOINT ["/app/main"]


