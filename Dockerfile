### Multistage Go build

# Build from source
FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /pds ./cmd/pds

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

COPY --from=build-stage /pds /pds

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/pds"]