FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/messages-svc ./cmd

FROM alpine:3.22

COPY --from=build /out/messages-svc /usr/local/bin/messages-svc

EXPOSE 9990

ENTRYPOINT ["messages-svc"]
