# Build a static binary
FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /shark .

# Minimal runtime image
FROM gcr.io/distroless/static-debian12
COPY --from=build /shark /shark
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/shark"]
