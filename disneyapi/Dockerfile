FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /disneyapi

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /disneyapi /disneyapi

EXPOSE 3000

USER nonroot

ENTRYPOINT ["/disneyapi"]