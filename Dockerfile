FROM golang:1.22.1 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY pkg/ ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prtglaw ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build /app/prtglaw ./prtglaw
EXPOSE 8888

ENTRYPOINT [ "/app/prtglaw" ]