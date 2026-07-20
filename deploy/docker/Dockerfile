FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/meclaw/meclaw/cmd.Version=docker" -o /out/meclaw .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/meclaw /usr/local/bin/meclaw
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["meclaw"]
CMD ["serve", "-c", "/etc/meclaw/config.json"]
