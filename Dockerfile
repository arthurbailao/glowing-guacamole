FROM golang:1.25.3 AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /cmd -mod=readonly ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /cmd /cmd

ENTRYPOINT ["/cmd"]
